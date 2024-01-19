package main

import (
	"context"
	"fmt"
	"time"
)

type Ci struct {
	Ctr         *Container
	Dir         *Directory
	dockerToken Optional[*Secret]
}

func (m *Ci) initBaseImage() {
	if m.Ctr == nil {
		m.Ctr = dag.Container().From("golang:alpine").
			WithMountedCache("/go/pkg/mod", dag.CacheVolume("snippetbox-go-mod")).
			WithMountedCache("/go/build-cache", dag.CacheVolume("snippetbox-go-build")).
			WithEnvVariable("GOCACHE", "/go/build-cache").
			WithExec([]string{"apk", "add", "tree"}).
			WithExec([]string{"apk", "add", "mysql-client"})
	}
}

// Run entire CI pipeline
// example usage: "dagger call -m ci ci --dir ."
func (m *Ci) Ci(ctx context.Context, dir *Directory) string {
	m.initBaseImage()

	ci := &Ci{
		Ctr: m.Ctr,
		Dir: dir,
	}

	output, _ := ci.Ctr.
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"echo", "it works!"}).
		Stdout(ctx)

	return output
}

// publish to dockerhub
func (m *Ci) Publish(
	ctx context.Context,
	dir *Directory,
	token Optional[*Secret],
	commit Optional[string],
) (string, error) {
	m.initBaseImage()

	ci := &Ci{
		Ctr: m.Ctr,
		Dir: dir,
	}

	dockerToken, isset := token.Get()
	gitCommit := commit.GetOr("latest")

	if isset {

		tags := [2]string{"latest", gitCommit}

		ctr := ci.Ctr.
			WithDirectory("/src", ci.Dir).
			WithRegistryAuth("docker.io", "levlaz", dockerToken)

		for _, tag := range tags {
			addr, err := ctr.Publish(ctx, fmt.Sprintf("levlaz/snippetbox:%s", tag))
			if err != nil {
				panic(err)
			}
			fmt.Println("Published: ", addr)
		}

		return fmt.Sprintf("Published: %s", tags), nil
	}

	return "Must pass registry token to publish", nil
}

// Serve development site
// example usage: "dagger up -m ci --port 4000:4000 serve --dir=."
func (m *Ci) Serve(dir *Directory) *Service {
	m.initBaseImage()

	ci := &Ci{
		Ctr: m.Ctr,
		Dir: dir,
	}

	mariadb := dag.Container().
		From("mariadb:latest").
		WithEnvVariable("MARIADB_ALLOW_EMPTY_ROOT_PASSWORD", "1").
		WithExposedPort(3306).
		AsService()

	return ci.Ctr.
		WithServiceBinding("db", mariadb).
		WithDirectory("/src", ci.Dir).
		WithWorkdir("/src").
		WithExposedPort(4000).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"sh", "-c", "mysql -h db -u root < internal/db/init.sql"}).
		WithExec([]string{"go", "run", "./cmd/web", "--dsn", "web:pass@tcp(db)/snippetbox?parseTime=true"}).
		AsService()
}
