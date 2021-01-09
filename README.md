# docker-thing

Command for docker things.

## Installation

TBA

## Features

```
Usage:
  dt now [flags]

Flags:
      --entry         Use the default entrypoint. If false you must provide one (default true)
  -h, --help          help for now
      --home          Mount the home directory of the user (default true)
      --pull          Pull the docker image even if present
      --pwd           Mount the PWD into the container (and set as working directory /pwd) (default true)
      --user string   User override for the command (default is current user) (default "CURRENTUSER")

Global Flags:
  -v, --verbose   verbose output
```

## Example usage

Run an interactive shell in the latest php image:

```sh
docker-thing now php -- -a
```

Run an bash in the latest ubuntu image (overriding default point):

```sh
docker-thing now --entry=0 --user=root ubuntu bash
```

Run composer version 1 info in the current working directory:

```sh
docker-thing now --pwd composer:1 info
```

Run git in the current working directory, with your home directory mounted

```sh
docker-thing now --home git config -- --list
```

Run a command with the tools verbose mode turned on

```sh
docker-thing now --verbose <image> <command>
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

Built using https://github.com/mitchellh/gox and a `Makefile`

```sh
go get github.com/mitchellh/gox
```

Then build...

```sh
make build
```

The `build` directory will be populated with the release.
