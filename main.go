package main

import (
	"github.com/addshore/dockerit/cmd"
)

var (
	// VERSION is set during build
	VERSION = "dev"
)

func main() {
	cmd.Execute(VERSION)
}
