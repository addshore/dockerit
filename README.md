# dockerit

[![Release](https://img.shields.io/github/release/addshore/dockerit.svg?style=flat-square)](https://github.com/addshore/dockerit/releases/latest)
![Tested Docker Versions](https://img.shields.io/badge/tested%20docker%20versions-18%2019%2020-blue)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/addshore/dockerit?style=flat-square)](https://goreportcard.com/report/github.com/addshore/dockerit)

Command for docker things.

## Installation

Head to the [releases page](https://github.com/addshore/dockerit/releases) and download the latest version for your system.

## Features

Containers are deleted after the command has run.

```
Usage:
  dockerit [image]
  dockerit [image] [command]
  dockerit [flags] [image] [command] -- [command flags]

Flags:
      --entry         Use the default entrypoint. If entry=0 you must provide one (default true)
  -h, --help          help for dockerit
      --home          Mount the home directory of the user
      --me            User override for the command, runs as current user
      --port string   Port mapping <host>:<container> eg. 8080:80
      --pull          Pull the docker image even if present
      --pwd           Mount the PWD into the container (and set as working directory /pwd)
      --user string   User override for the command
  -v, --verbose       verbose output
      --version       version infomation
```

## Example usage

Output help infomation:

```sh
dockerit --help
```

Run an interactive shell in the latest php image:

```sh
dockerit php -- -a
```

Run an bash in the latest ubuntu image (overriding default point):

```sh
dockerit --entry=0 --user=root ubuntu bash
```

Run composer version 1 info in the current working directory as the current user with their home dir mounted:

```sh
dockerit --me --pwd --home composer:1 info
```

Run nginx as the container user and expose it on port 8080:

```sh
dockerit --port=8080:80 nginx
```

Run git in the current working directory as the current user with their home dir mounted:

```sh
dockerit --me --pwd --home git config -- --list
```

Run a command in an image with verbose mode turned on:

```sh
dockerit --verbose [image] [command]
```


## Development

### Dependencies

```sh
go mod vendor
```

### Run from source

There is a shortcut to help you:

```sh
./dev
```

### Building

Built using https://github.com/laher/goxc and a `Makefile`

```sh
go get github.com/laher/goxc
```

Then build...

```sh
make build
```

If you want to build with a specific version number / name:

```sh
VERSION=0.0.0 make build
```

The `build` directory will be populated with the release.
