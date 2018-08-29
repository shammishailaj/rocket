# Introduction

## Installation


### Latest
```bash
curl -sSf https://raw.githubusercontent.com/astrocorp42/rocket/master/install.sh | sh
```

### Binary releases
[https://github.com/astrocorp42/rocket/releases/latest](https://github.com/astrocorp42/rocket/releases/latest)

### Using go (nightly)
```bash
$ go get -u github.com/astrocorp42/rocket
```


## CI usage

You may want to use the `ci.sh` script to ease the usage in a CI/CD environment.
Here an example with travis
```yml
# .travis.yml

sudo: false
language: go
env:
  - GO111MODULE=on

go:
  - 1.11

script:
  - make

deploy:
  provider: script
  skip_cleanup: true
  # The important part: it's the same as
  # curl -sSf https://raw.githubusercontent.com/astrocorp42/rocket/master/install.sh && $HOME/.rocket/rocket
  # you can pass argument: ...ci.sh | sh -s -- -c abc/xys/another_config_file.toml
  script: curl -sSf https://raw.githubusercontent.com/astrocorp42/rocket/master/ci.sh | sh
  on:
    repo: astrocorp42/rocket
    tags: true
```

## Environment variables
