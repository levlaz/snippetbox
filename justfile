# show all available options
default:
    just --list 

# run ci pipleine
ci:
    dagger call -m ci ci --dir .

# serve site on localhost:4000
serve:
    dagger up --focus=false -m ci --port 4000:4000 serve --dir=.

# run ci with empty commit 
empty: 
    git commit --allow-empty -m 'trigger ci' && git push

# push to all remotes
push: 
    git push -u gitlab main && git push -u origin main

# publish docker image to dockerhub 
publish:
    dagger call -m ci publish --dir . --token env:DOCKER_TOKEN --commit $(git rev-parse HEAD)

# start dev database 
db:
    dagger up -m ci/mariadb --port 3306:3306 serve

# init database 
init-db: 
    mysql -h 127.0.0.1 -u root < internal/db/init.sql

# you can pass this in yourself too by running
# just mysql < internal/db/seed.sql
# seed db
seed-db:
    mysql -h 127.0.0.1 -u root snippetbox < internal/db/seed.sql

# start mysql-client with snippetbox database 
mysql: 
    mysql -h 127.0.0.1 -u root snippetbox