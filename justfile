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

# push to gitlab
push: 
    git push -u gitlab main && git push -u origin main