package main

import (
	"fmt"
	"os"
	"strings"
)

type Command struct {
	args []string
}

type Options struct {
	Help             string
	PathNotExist     string
	NotEnoughArgs    string
	OutputFileNeeded string
	Version          string
}

func (c *Command) Check() error {
	options := Options{
		Help: `gqlmerge is the tool to merge & stitch *.graphql files and generate a Graphql schema

	Usage:	gqlmerge [PATH] [OUTPUT.graphql]

	e.g.

		gqlmerge ./schema schema.graphql

	Options:

		-v	: check the version
		-h	: help
`,
		PathNotExist:     "Path '%s' does not Exist",
		NotEnoughArgs:    "Not enough arguments",
		OutputFileNeeded: "Output file argument is needed",
		Version:          "v0.1.0",
	}
	// show the version
	if strings.HasPrefix(c.args[1], "-v") {
		return fmt.Errorf(options.Version)
	} else if strings.HasPrefix(c.args[1], "-h") {
		// show the version
		return fmt.Errorf(options.Help)
	}

	// check the number of args
	if len(c.args) <= 1 {
		return fmt.Errorf(options.NotEnoughArgs)
	} else if len(c.args) == 2 {
		return fmt.Errorf(options.OutputFileNeeded)
	}

	// check first arg, path is existing
	if _, err := os.Stat(c.args[1]); os.IsNotExist(err) {
		return fmt.Errorf(options.PathNotExist, c.args[1])
	}

	return nil
}
