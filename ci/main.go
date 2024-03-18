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
func (m *Ci) Lint(ctx context.Context, dir *Directory) (string, error) {
	return dag.Golang().
		WithProject(dir).
		GolangciLint(ctx)
}

// Run test suite
func (m *Ci) Test(ctx context.Context, dir *Directory) (string, error) {
	return m.base().
		WithDirectory("/src", dir).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "./cmd/web"}).
		Stdout(ctx)
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
) string {

	var output string

	// run linter
	lintOutput, err := m.Lint(ctx, dir)
	if err != nil {
		fmt.Sprint(err)
	}
	output = output + "\n" + lintOutput

	// run tests
	testOutput, err := m.Test(ctx, dir)
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

// Debug Build Container
func (m *Ci) Debug(dir *Directory) *Container {
	return m.base().
		WithServiceBinding("db", dag.Mariadb().Serve()).
		WithDirectory("/src", dir).
		WithWorkdir("/src")
}
