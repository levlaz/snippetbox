# show all available options
default:
    just --list 

# run ci pipleine
ci:
    dagger call -m ci ci --dir .

# serve site on localhost:4000
serve:
    dagger up -m ci --port 4000:4000 serve --dir=.