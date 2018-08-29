# Introduction

`rocket` is a deployment automation tool. It's the `D` in`CI/CD`. The goal is to provide an easy
to use uniform experience of software deployement whether in a CI environment or on your local laptotp.

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

## Usage

Start by creating a `.rocket.toml`


## Environments

`rocket` support different environments throught different configuration files:
```
$ tree -a
.
├── .rocket_dev.toml
└── .rocket.toml
```
then you ccan run
```bash
$ rocket # -> use the default .rocket.toml
$ rocket -c .rocket_dev.toml # to deploy in your dev environment
```

## CI usage

You may want to use the `ci.sh` script to ease the usage in a CI/CD environment to ease `rocket` installation and usage.
Here an example with travis
```yaml
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

When starting **rocket** prepares the deploy environment. It starts by setting a list of **predefined environment variables** and a list of **user-defined environment variables**.

### Priority of variables

The variables can be overwritten and they take precedence over each other in this order:

1. Already set environment variables (take precedence over all)
2. [TOML-defined environment variables](#toml-defined-environment-variables)
3. [Predefined variables](#predefined-environment-variables) (are the lowest in the chain)

### TOML-defined environment variables

rocket allow you to define variables inside `.rocket.toml` that are then injected in the environment.
For example:
```toml
[env]
MY_VARIABLE = "MYSUPERVALUE"
# You are able to use other variables inside your variable definition (or escape them with $$):
FLAGS = "-al"
LS_CMD = "ls $FLAGS $$TMP_DIR" # -> 'ls -al $TMP_DIR'

[heroku]
api_key = $HEROKU_TOKEN # -> it's not defined above nor in the predefined variables, so it will expand to the already set environment variable
```

### Predefined environment variables

| Variable             | Description |
| --------------------- | -------|
| **ROCKET_COMMIT_HASH** | The current commit revision |
| **ROCKET_LAST_TAG** | The last commit tag name |
| **ROCKET_GIT_REPO** |  The slug (in form: **owner_name/repo_name**) of the repository currently being deployed |

## Roadmap

See [https://github.com/astrocorp42/rocket/projects/2](https://github.com/astrocorp42/rocket/projects/2)
