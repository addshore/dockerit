package main

import (
	"strings"

	"github.com/addshore/dockerit/cmd"
)

var (
	// VERSION is set during build
	// Surrounded by ' as this apparently happens? https://github.com/addshore/dockerit/issues/9
	// They are trimmed off further down...
	VERSION = "'dev'"
	// SOURCE_DATE is set during build
	SOURCE_DATE = "unknown"
)

func main() {
	cmd.Execute(strings.Trim(VERSION,"'"), SOURCE_DATE)
}
