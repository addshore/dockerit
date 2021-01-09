# docker-thing

Command for docker things.

## Guide

```sh
docker-thing now composer:1 help
```

### Features

TBA

## Development

### Installation

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
