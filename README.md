# Rocket

[![GoDoc](https://godoc.org/github.com/astrocorp42/rocket?status.svg)](https://godoc.org/github.com/astrocorp42/rocket)
[![GitHub release](https://img.shields.io/github/release/astrocorp42/rocket.svg)](https://github.com/astrocorp42/rocket/releases/latest)
[![Build Status](https://travis-ci.org/astrocorp42/rocket.svg?branch=master)](https://travis-ci.org/astrocorp42/rocket)

Deploy software as fast and easily as possible. The `D` in CI/CD.

1. [Installation](#installation)
2. [Usage](#usage)
3. [Documentation](#documentation)
4. [Available providers](#available-providers)
5. [Roadmap](#roadmap)

-------------------


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

Go to your project's root directory then
```bash
$ rocket init # create a configuration .rocket.toml file with default configuration
# edit the file with the desired configuration
$ cat .rocket.toml
```
```toml
description = "This is a configuration file for rocket: Deploy software as fast and easily as possible. See https://github.com/astrocorp42/rocket"

[github_releases]
assets = [
  "dist/*.zip",
  "dist/rocket_*_sha512sums.txt"
]
```
```bash
$ rocket # to deploy
```



## Documentation
See [https://astrocorp.net/rocket](https://astrocorp.net/rocket)




## Available providers

| Provider              | Status | Documentation |
| --------------------- | -------| ------------- |
| [AWS S3](https://aws.amazon.com/s3) `s3` | :construction: | - |
| Custom script `script` | :heavy_check_mark: | [docs](https://astrocorp.net/rocket/custom_script) |
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
