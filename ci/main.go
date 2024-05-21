package main

import (
	"context"
	"fmt"
	"time"
)

type Ci struct{}

func (m *Ci) base() *Container {
	return dag.Container().From("golang:alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("snippetbox-go-mod")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("snippetbox-go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"apk", "add", "tree"}).
		WithExec([]string{"apk", "add", "mysql-client"})
}

// Lint
func (m *Ci) Lint(ctx context.Context, dir *Directory) *Container {
 	return dag.GolangciLint().Run(dir)
}

// Run test suite
func (m *Ci) Test(ctx context.Context, dir *Directory) *Container {
	return m.base().
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "./cmd/web"})
}

// Run entire CI pipeline
// example usage: "dagger call -m ci ci --dir ."
func (m *Ci) Ci(
	ctx context.Context,
	dir *Directory,
	// +optional
	token *Secret,
	// +optional
	// +default="latest"
	commit string,
	// +optional
	// +default="local"
	// environment where ci build is running, workaround for dagger cloud not
	// showing this information yet.
	env string,
) string {

	var output string

	// run linter
	lintOutput, err := m.Lint(ctx, dir).Stdout(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + lintOutput

	// run tests
	testOutput, err := m.Test(ctx, dir).Stdout(ctx)
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

	// get TraceURL
	traceUrl, err := dag.CloudUtils().TraceURL(ctx)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n\nDagger Cloud Trace: " + traceUrl

	return output
}

// publish to dockerhub
func (m *Ci) Publish(
	ctx context.Context,
	dir *Directory,
	// +optional
	token *Secret,
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
	}

	return "Must pass registry token to publish", nil
}

// Serve development site
// example usage: "dagger call serve --dir=. up."
func (m *Ci) Serve(dir *Directory) *Service {
	return m.base().
		WithServiceBinding("db", dag.Mariadb().Serve()).
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExposedPort(4000).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "mysql -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"sh", "-c", "mysql -h db -u root snippetbox < internal/db/seed.sql"}).
		WithExec([]string{"go", "run", "./cmd/web", "--dsn", "web:pass@tcp(db)/snippetbox?parseTime=true"}).
		AsService()
}

// Debug build container with MariaDB service attached
func (m *Ci) Debug(dir *Directory) *Container {
	return m.base().
		WithServiceBinding("db", dag.Mariadb().Serve()).
        WithServiceBinding("dragonfly", dag.Dragonfly().Serve()).
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-c", "mysql -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"sh", "-c", "mysql -h db -u root snippetbox < internal/db/seed.sql"})
}

// Get Private Container
func (m *Ci) WithPrivateContainer(
	// full private image address in the form of [host]/[user]/[repo]:[tag]
	address string,
	// registry username
	username string,
	// registry token
	token *Secret,
) *Container {
	return dag.Container().WithRegistryAuth(address, username, token)
}
