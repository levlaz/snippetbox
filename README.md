# Snippetbox 

## Connecting to DB Locally 

```bash
# start db service with dagger 
dagger up -m ci/mariadb --port 3306:3306 serve

# connect with mysql-client
mysql -h 127.0.0.1 -u root

# run init stuff
mysql -h 127.0.0.1 -u root < internal/db/init.sql
```

## CI 

This project intentionally runs CI on every known SaaS provider using Dagger.

* CircleCI -> `.circleci/config.yml`
* GitHub Actions -> `.github/workflows/dagger.yml`
* GitLab -> `.gitlab-ci.yml`
* SourceHut -> `.build.yml`