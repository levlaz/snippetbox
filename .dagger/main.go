package main

import (
	"context"
	"fmt"
	"main/internal/dagger"
	"time"
)

func New(
	// +defaultPath="/"
	src *dagger.Directory,
) *Snippetbox {
	return &Snippetbox{
		Src: src,
	}
}

type Snippetbox struct {
	Src *dagger.Directory
}

func (m *Snippetbox) base() *dagger.Container {
	return dag.Container().From("golang:alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("snippetbox-go-mod")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("snippetbox-go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"apk", "add", "tree"}).
		WithExec([]string{"apk", "add", "mysql-client"}).
		WithExec([]string{"apk", "add", "golangci-lint"})
}

// Lint
func (m *Snippetbox) Lint(
	ctx context.Context,
) *dagger.Container {
	return m.base().
		WithDirectory("/src", m.Src).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v"})
}

// Build snippetbox binary for all supported platforms
func (m *Snippetbox) Build(
	ctx context.Context,
) *dagger.Directory {
	// define build matrix
	gooses := []string{"linux", "darwin", "windows"}
	goarches := []string{"amd64", "arm64"}

	// create empty directory to put build artifacts
	outputs := dag.Directory()

	// run build for each combination
	for _, goos := range gooses {
		for _, goarch := range goarches {
			// create directory for each OS and architecture
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// build artifact
			build := m.base().
				WithDirectory("/src", m.Src).
				WithWorkdir("/src").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", path, "./cmd/web"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	return outputs
}

// Run test suite
func (m *Snippetbox) Test(
	ctx context.Context,
	// quiet output for tests
	// +optional
	// +default=false
	quiet bool,
) *dagger.Container {
	ctr := m.base().
		WithDirectory("/src", m.Src).
		WithWorkdir("/src").
		WithEnvVariable("CACHEBUSTER", time.Now().String())

	if quiet {
		ctr = ctr.WithExec([]string{"go", "test", "./..."})
	} else {
		ctr = ctr.WithExec([]string{"go", "test", "-v", "./..."})
	}

	return ctr
}

// publish to dockerhub or ttl.sh if no token is provided
func (m *Snippetbox) Publish(
	ctx context.Context,
	// +optional
	token *dagger.Secret,
	// +optional
	// +default="latest"
	commit string,
) (string, error) {
	if token != nil {
		ctr := m.base().
			WithDirectory("/src", m.Src).
			WithRegistryAuth("docker.io", "levlaz", token)

		addr, err := ctr.Publish(ctx, fmt.Sprintf("levlaz/snippetbox:%s", commit))
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}

		return fmt.Sprintf("Published: %s", addr), nil
	} else {
		addr, err := m.base().
			WithDirectory("/src", m.Src).
			Publish(ctx, fmt.Sprintf("ttl.sh/levlaz/snippetbox:%s", commit))
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}

		return fmt.Sprintf("Published: %s", addr), nil
	}
}

// Return snippetbox server with database service attached
// example usage: "dagger call server up"
func (m *Snippetbox) Server(
	// +optional
	database *dagger.Service,
) *dagger.Container {
	if database == nil {
		database = dag.Mariadb().Serve()
	}
	return m.base().
		WithServiceBinding("db", database).
		WithDirectory("/src", m.Src).
		WithWorkdir("/src").
		WithExposedPort(4000).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "mariadb --skip-ssl -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"sh", "-c", "mariadb --skip-ssl -h db -u root snippetbox < internal/db/seed.sql"}).
		WithDefaultArgs([]string{"go", "run", "./cmd/web", "--dsn", "web:pass@tcp(db)/snippetbox?parseTime=true"})
}

// Run entire CI pipeline
// example usage: "dagger call ci"
func (m *Snippetbox) Ci(
	ctx context.Context,
	// +optional
	token *dagger.Secret,
	// +optional
	// +default="latest"
	commit string,
) string {

	var output string

	// run linter
	lintOutput, err := m.Lint(ctx).Stdout(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + lintOutput

	// run tests
	testOutput, err := m.Test(ctx, false).Stdout(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + testOutput

	// publish container
	publishOutput, err := m.Publish(ctx, token, commit)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + publishOutput

	return output
}

// return container with service attached that is not running
func (m *Snippetbox) Debug(
	// +optional
	database *dagger.Service,
) *dagger.Container {
	if database == nil {
		database = dag.Mariadb().Serve()
	}
	return m.base().
		WithServiceBinding("db", database).
		WithDirectory("/src", m.Src).
		WithWorkdir("/src")
}

// Return private container
func (m *Snippetbox) WithPrivateContainer(
	// full private image uri
	uri string,
	// registry uri
	address string,
	// registry username
	username string,
	// registry password
	token *dagger.Secret,
) *dagger.Container {
	return dag.
		Container().
		WithRegistryAuth(address, username, token).
		From(uri)
}

// Publish and reset auth
func (m *Snippetbox) PublishAndResetAuth(
	ctx context.Context,
	// full private image uri
	uri string,
	// registry uri
	address string,
	// registry username
	username string,
	// registry password
	token *dagger.Secret,
	// directory with Dockerfile
	src *dagger.Directory,
) (string, error) {
	// auth with registry
	dag.Container().WithRegistryAuth(address, username, token).Sync(ctx)

	// build and publish
	out, err := src.DockerBuild().Publish(ctx, uri)
	if err != nil {
		return "", fmt.Errorf("failed to publish: %w", err)
	}

	// try to pull from private image (should succeed)
	workingPrivateOut, err := dag.Container().From(uri).WithExec([]string{"echo", "hello from working private image"}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to pull private image: %w", err)
	}

	out += workingPrivateOut

	// reset auth
	dag.Container().WithRegistryAuth(address, "username", dag.SetSecret("nonsense", "nonsense")).Sync(ctx)

	// try to pull from private image (should fail)
	nonWorkingPrivateOut, err := dag.Container().From(uri).WithExec([]string{"echo", "hello from nonworking private image"}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to pull private image: %w", err)
	}
	out += nonWorkingPrivateOut

	return out, nil
}

// Publish and reset auth test using push
func (m *Snippetbox) PublishAndResetAuthTest(
	ctx context.Context,
	// full private image uri
	uri string,
	// registry uri
	address string,
	// registry username
	username string,
	// registry password
	token *dagger.Secret,
	// directory with Dockerfile
	src *dagger.Directory,
) (string, error) {
	// auth with registry
	dag.Container().WithRegistryAuth(address, username, token).Sync(ctx)

	// build and publish with creds
	out, err := src.DockerBuild().Publish(ctx, uri)
	if err != nil {
		return "", fmt.Errorf("failed to publish: %w", err)
	}

	// does not seem to actually reset auth, old auth is being cached and resued
	dag.Container().WithoutRegistryAuth(address).Sync(ctx)

	// dag.Container().WithRegistryAuth(address, username, dag.SetSecret("token", "reset")).Sync(ctx)

	// build and publish without creds
	withOut, err := src.WithNewFile("hi.txt", "hello").DockerBuild().Publish(ctx, uri)
	if err != nil {
		return "", fmt.Errorf("failed to publish: %w", err)
	}

	out += withOut

	return out, nil
}
