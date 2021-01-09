package main

import (
	"github.com/addshore/docker-thing/cmd"
)

var (
	// VERSION is set during build
	VERSION = "0.0.0"
)

func main() {
	cmd.Execute()
}
