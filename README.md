# Rocket

[![GoDoc](https://godoc.org/github.com/astrocorp42/rocket?status.svg)](https://godoc.org/github.com/astrocorp42/rocket)
[![GitHub release](https://img.shields.io/github/release/astrocorp42/rocket.svg)](https://github.com/astrocorp42/rocket/releases/latest)
[![Build Status](https://travis-ci.org/astrocorp42/rocket.svg?branch=master)](https://travis-ci.org/astrocorp42/rocket)
[![Docker image](https://img.shields.io/badge/docker-astrocorp/rocket-blue.svg)](https://hub.docker.com/r/astrocorp/rocket)


Automated deployment as fast and easy as possible. `rocket` is the `D` in CI/CD: It allows to easily deploy software across a large range of providers from any CI/CD pipeline.

1. [Installation](#installation)
2. [Docker image](#docker-image)
3. [Usage](#usage)
4. [Documentation](#documentation)
5. [Available providers](#available-providers)
6. [Roadmap](#roadmap)

-------------------


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

`astrocorp/rocket`


## Usage

Go to your project's root directory then
```bash
$ rocket init # create a configuration .rocket.toml file with default configuration
# edit the file with the desired configuration
$ cat .rocket.toml
```
```toml
description = "This is a configuration file for rocket: Automated deployment as fast and easy as possible. See https://github.com/astrocorp42/rocket"

[github_releases]
assets = [
  "dist/*.zip",
  "dist/rocket_*_sha512sums.txt"
]
```
```bash
$ rocket # to deploy
```

See [https://github.com/astrocorp42/rocket/blob/master/.rocket.toml](https://github.com/astrocorp42/rocket/blob/master/.rocket.toml) for an example with the `github_releases` provider.

Help
```bash
$ rocket help
Automated deployment as fast and easy as possible. rocket is the D in CI/CD. See https://github.com/astrocorp42/rocket

Usage:
  rocket [flags]
  rocket [command]

Available Commands:
  help        Help about any command
  init        Init rocket by creating a .rocket.toml configuration file
  version     Display the version and build information

Flags:
  -c, --config string   Use the specified configuration file (and set it's directory as the working directory
  -d, --debug           Display debug information
  -h, --help            help for rocket

Use "rocket [command] --help" for more information about a command.
```

## Documentation

See [https://astrocorp.net/rocket](https://astrocorp.net/rocket)




## Available providers

| Provider              | Status | Documentation |
| --------------------- | -------| ------------- |
| [AWS S3](https://aws.amazon.com/s3) `s3` | :construction: | - |
| Custom script `script` | :heavy_check_mark: | [docs](https://astrocorp.net/rocket/custom_script) |
| [Docker](https://www.docker.com) `docker` | :construction: | - |
| [Google Firebase](https://firebase.google.com) `firebase` | :clock1: | - |
| [Google Cloud Storage](https://cloud.google.com/storage) `gcs` | :clock1: | - |
| [GitHub releases](https://help.github.com/categories/releases) `github_releases` | :heavy_check_mark: | [docs](https://astrocorp.net/rocket/github_releases) |
| [Heroku](https://www.heroku.com) `heroku` | :heavy_check_mark: | [docs](https://astrocorp.net/rocket/heroku) |
| [Netlify](https://www.netlify.com) `netlify` | :clock1: | - |
| [NPM](https://www.npmjs.com) `npm` | :clock1: | - |
| [SCP](https://en.wikipedia.org/wiki/Secure_copy) `scp` | :clock1: | - |
| [SFTP](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol) `sftp` | :clock1: | - |
| [SSH](https://en.wikipedia.org/wiki/Secure_Shell) `ssh` | :clock1: | - |

:heavy_check_mark: = Done :construction: = in progress :clock1: = planned


## Roadmap

See [https://github.com/astrocorp42/rocket/projects/2](https://github.com/astrocorp42/rocket/projects/2)
