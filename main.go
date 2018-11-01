package main

import (
	"github.com/eschizoid/flixctl/cmd"
)

var (
	BUILD   string
	VERSION string
)

func main() {
	cmd.Execute(VERSION, BUILD)
}
