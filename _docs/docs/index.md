# Introduction

`rocket` is a software delivery automation tool. It's the `D` in`CI/CD`. The goal is to provide an easy
to use uniform experience of software delivery whether in a CI environment or on your local laptotp.

It allows to easily release software across a large range of providers from any CI/CD pipeline.

The only required dependecy is `git`.


[![GoDoc](https://godoc.org/github.com/astrocorp42/rocket?status.svg)](https://godoc.org/github.com/astrocorp42/rocket)
[![GitHub release](https://img.shields.io/github/release/astrocorp42/rocket.svg)](https://github.com/astrocorp42/rocket/releases/latest)
[![Build Status](https://travis-ci.org/astrocorp42/rocket.svg?branch=master)](https://travis-ci.org/astrocorp42/rocket)
[![Docker image](https://img.shields.io/badge/docker-astrocorp/rocket-blue.svg)](https://hub.docker.com/r/astrocorp/rocket)



## Installation

### Using go (nightly)
```bash
$ go get -u github.com/astrocorp42/rocket
```

### Latest
```bash
curl -sSf https://raw.githubusercontent.com/astrocorp42/rocket/master/install.sh | sh
```

### Binary releases
[https://github.com/astrocorp42/rocket/releases/latest](https://github.com/astrocorp42/rocket/releases/latest)




## Docker image

[astrocorp/rocket](https://hub.docker.com/r/astrocorp/rocket)



## Available providers

| Provider              | Status | Documentation |
| --------------------- | -------| ------------- |
| [AWS S3](https://aws.amazon.com/s3) `aws_s3` | 🚧 | [docs](https://astrocorp.net/rocket/aws_s3) |
| Custom script `script` | ✅| [docs](https://astrocorp.net/rocket/custom_script) |
| [Docker](https://www.docker.com) `docker` | ✔ | [docs](https://astrocorp.net/rocket/docker) |
| [Google Firebase](https://firebase.google.com) `firebase` | 🕐 | - |
| [Google Cloud Storage](https://cloud.google.com/storage) `gcs` | 🕐 | - |
| [GitHub releases](https://help.github.com/categories/releases) `github_releases` | ✔ | [docs](https://astrocorp.net/rocket/github_releases) |
| [Heroku](https://www.heroku.com) `heroku` | ✔ | [docs](https://astrocorp.net/rocket/heroku) |
| [Netlify](https://www.netlify.com) `netlify` | 🕐 | - |
| [NPM](https://www.npmjs.com) `npm` | 🕐 | - |
| [SCP](https://en.wikipedia.org/wiki/Secure_copy) `scp` | 🕐 | - |
| [SFTP](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol) `sftp` | 🕐 | - |
| [SSH](https://en.wikipedia.org/wiki/Secure_Shell) `ssh` | 🕐 | - |

✔ = Done 🚧 = in progress 🕐 = planned




## Usage

See [https://github.com/astrocorp42/rocket/blob/master/.rocket.toml](https://github.com/astrocorp42/rocket/blob/master/.rocket.toml) for an example with the `github_releases` provider.

Start by creating a `.rocket.toml` file. Here is an example to deploy a GitHub release:
```toml
description = "This is a configuration file for rocket: automated software delivery as fast and easy as possible. See https://github.com/astrocorp42/rocket"

[github_releases]
# the assets to upload
assets = [
  "dist/*.zip",
  "dist/rocket_*_sha512sums.txt"
]
```



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
api_key = "$HEROKU_TOKEN" # -> it's not defined above nor in the predefined variables, so it will expand to the already set environment variable
```

### Predefined environment variables

| Variable             | Description |
| --------------------- | -------|
| **ROCKET_COMMIT_HASH** | The current commit revision |
| **ROCKET_LAST_TAG** | The last commit tag name |
| **ROCKET_GIT_REPO** |  The slug (in form: **owner_name/repo_name**) of the repository currently being deployed |




## Roadmap

See [https://github.com/astrocorp42/rocket/projects/2](https://github.com/astrocorp42/rocket/projects/2)
