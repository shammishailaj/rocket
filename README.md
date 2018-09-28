<p align="center">
  <img alt="rocket logo" src="https://astrocorp.net/imgs/landing/rocket.png" height="180" />
  <h3 align="center">Rocket</h3>
  <p align="center">Automated software delivery as fast and easy as possible</p>
</p>

--------

[![GoDoc](https://godoc.org/github.com/astrocorp42/rocket?status.svg)](https://godoc.org/github.com/astrocorp42/rocket)
[![GitHub release](https://img.shields.io/github/release/astrocorp42/rocket.svg)](https://github.com/astrocorp42/rocket/releases/latest)
[![Build Status](https://travis-ci.org/astrocorp42/rocket.svg?branch=master)](https://travis-ci.org/astrocorp42/rocket)
[![Docker image](https://img.shields.io/badge/docker-astrocorp/rocket-blue.svg)](https://hub.docker.com/r/astrocorp/rocket)

`rocket` is the `D` in CI/CD: It allows to easily deliver software to a large range of providers from any CI/CD pipeline.

1. [Installation](#installation)
2. [Docker image](#docker-image)
3. [Available providers](#available-providers)
4. [Usage](#usage)
5. [Documentation](#documentation)
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

[astrocorp/rocket](https://hub.docker.com/r/astrocorp/rocket)




## Available providers

| Provider              | Status | Documentation |
| --------------------- | -------| ------------- |
| [AWS Elastic Beanstalk](https://aws.amazon.com/elasticbeanstalk/) `aws_eb` | ✔ | [docs](https://astrocorp.net/rocket/aws_eb) |
| [AWS S3](https://aws.amazon.com/s3) `aws_s3` | ✔ | [docs](https://astrocorp.net/rocket/aws_s3) |
| Custom script `script` | ✔ | [docs](https://astrocorp.net/rocket/custom_script) |
| [Docker](https://www.docker.com) `docker` | ✔ | [docs](https://astrocorp.net/rocket/docker) |
| [Google Firebase](https://firebase.google.com) `firebase` | 🕐 | - |
| [Google Cloud Storage](https://cloud.google.com/storage) `gcs` | 🕐 | - |
| [GitHub releases](https://help.github.com/categories/releases) `github_releases` | ✔ | [docs](https://astrocorp.net/rocket/github_releases) |
| [Heroku](https://www.heroku.com) `heroku` | ✔ | [docs](https://astrocorp.net/rocket/heroku) |
| [Netlify](https://www.netlify.com) `netlify` | 🚧 | - |
| [NPM](https://www.npmjs.com) `npm` | 🕐 | - |
| [SCP](https://en.wikipedia.org/wiki/Secure_copy) `scp` | 🕐 | - |
| [SFTP](https://en.wikipedia.org/wiki/SSH_File_Transfer_Protocol) `sftp` | 🕐 | - |
| [SSH](https://en.wikipedia.org/wiki/Secure_Shell) `ssh` | 🕐 | - |
| [ZEIT Now](https://zeit.co/now) `zeit_now` | ✔ | [docs](https://astrocorp.net/rocket/zeit_now) |

✔ = Done 🚧 = in progress 🕐 = planned





## Usage

`rocket` use [SAN](https://astrocorp.net/san) file for it's configuration.

Go to your project's root directory then
```bash
$ rocket init # create a configuration .rocket.san file with default configuration
# edit the file with the desired configuration
$ cat .rocket.san
```
```san
description = "This is a configuration file for rocket: Deploy software as fast and easily as possible. See https://github.com/astrocorp42/rocket"

github_releases = {
  assets = [
    "dist/*.zip",
    "dist/rocket_*_sha512sums.txt"
  ]
}

docker = {
  images = [
    "astrocorp/rocket:$VERSION",
    "astrocorp/rocket:latest"
  ]
}
```
```bash
$ rocket # to deploy
```

See [https://github.com/astrocorp42/rocket/blob/master/.rocket.san](https://github.com/astrocorp42/rocket/blob/master/.rocket.san) for an example using both the `github_releases` and the `docker` providers.



Help
```bash
$ rocket help
Automated software delivery as fast and easy as possible. rocket is the D in CI/CD. See https://github.com/astrocorp42/rocket

Usage:
  rocket [flags]
  rocket [command]

Available Commands:
  help        Help about any command
  init        Init rocket by creating a .rocket.san configuration file
  version     Display the version and build information

Flags:
  -c, --config string   Use the specified configuration file (and set it's directory as the working directory
  -d, --debug           Display debug information
  -h, --help            help for rocket

Use "rocket [command] --help" for more information about a command.
```




## Documentation

See [https://astrocorp.net/rocket](https://astrocorp.net/rocket)





## Roadmap

See [https://github.com/astrocorp42/rocket/projects/2](https://github.com/astrocorp42/rocket/projects/2)
