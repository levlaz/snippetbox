package main

import (
	"context"
	"fmt"
	"main/internal/dagger"
	"time"
)

type Snippetbox struct{}

func (m *Snippetbox) base() *dagger.Container {
	return dag.Container().From("golang:alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("snippetbox-go-mod")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("snippetbox-go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"apk", "add", "tree"}).
		WithExec([]string{"apk", "add", "mysql-client"})
}

// Lint
func (m *Snippetbox) Lint(
	ctx context.Context,
	// +defaultPath="/"
	dir *dagger.Directory,
) *dagger.Container {
	return dag.GolangciLint().Run(dir)
}

// Build snippetbox binary for all supported platforms
func (m *Snippetbox) Build(
	ctx context.Context,
	// +defaultPath="/"
	dir *dagger.Directory,
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
				WithDirectory("/src", dir).
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
	// +defaultPath="/"
	dir *dagger.Directory,
	// quiet output for tests
	// +optional
	// +default=false
	quiet bool,
) *dagger.Container {
	ctr := m.base().
		WithDirectory("/src", dir).
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
	// +defaultPath="/"
	dir *dagger.Directory,
	// +optional
	token *dagger.Secret,
	// +optional
	// +default="latest"
	commit string,
) (string, error) {
	if token != nil {
		ctr := m.base().
			WithDirectory("/src", dir).
			WithRegistryAuth("docker.io", "levlaz", token)

		addr, err := ctr.Publish(ctx, fmt.Sprintf("levlaz/snippetbox:%s", commit))
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}

		return fmt.Sprintf("Published: %s", addr), nil
	} else {
		addr, err := m.base().
			WithDirectory("/src", dir).
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
	// +defaultPath="/"
	dir *dagger.Directory,
	// +optional
	database *dagger.Service,
) *dagger.Container {
	if database == nil {
		database = dag.Mariadb().Serve()
	}
	return m.base().
		WithServiceBinding("db", database).
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExposedPort(4000).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "mysql -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"sh", "-c", "mysql -h db -u root snippetbox < internal/db/seed.sql"}).
		WithExec([]string{"go", "run", "./cmd/web", "--dsn", "web:pass@tcp(db)/snippetbox?parseTime=true"})
}

// Run entire CI pipeline
// example usage: "dagger call ci"
func (m *Snippetbox) Ci(
	ctx context.Context,
	// +defaultPath="/"
	dir *dagger.Directory,
	// +optional
	token *dagger.Secret,
	// +optional
	// +default="latest"
	commit string,
) string {

	var output string

	// run linter
	lintOutput, err := m.Lint(ctx, dir).Stdout(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + lintOutput

	// run tests
	testOutput, err := m.Test(ctx, dir, false).Stdout(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + testOutput

	// publish container
	publishOutput, err := m.Publish(ctx, dir, token, commit)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + publishOutput

	return output
}

// return container with service attached that is not running
func (m *Snippetbox) Debug(
	// +defaultPath="/"
	dir *dagger.Directory,
	// +optional
	database *dagger.Service,
) *dagger.Container {
	if database == nil {
		database = dag.Mariadb().Serve()
	}
	return m.base().
		WithServiceBinding("db", database).
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExposedPort(4000).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "mysql -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"sh", "-c", "mysql -h db -u root snippetbox < internal/db/seed.sql"})
}
