package lib

import (
	"fmt"
	// "os"
	"strings"
)

// Command for gqlmerge
type Command struct {
	Args []string
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
		Help: `👋 'gqlmerge' is the tool to merge & stitch *.graphql files and generate a Graphql schema
Author : Woonki Moon <woonki.moon@gmail.com>

Usage:	gqlmerge [PATH] [OUTPUT.graphql]

e.g.

	gqlmerge ./schema schema.graphql

Options:

	-v	: check the version
	-h	: help
`,
		PathNotExist:     "❌ Path '%s' does not Exist",
		NotEnoughArgs:    "❌ Not enough arguments",
		OutputFileNeeded: "❌ Output file argument is needed",
		Version:          "v0.1.1",
	}
	// show the version
	if strings.HasPrefix(c.Args[1], "-v") {
		return fmt.Errorf(options.Version)
	} else if strings.HasPrefix(c.Args[1], "-h") {
		// show the version
		return fmt.Errorf(options.Help)
	}

	// check the number of args
	if len(c.Args) <= 1 {
		return fmt.Errorf(options.NotEnoughArgs)
	} else if len(c.Args) == 2 {
		return fmt.Errorf(options.OutputFileNeeded)
	}

	// check first arg, path is existing
	// if _, err := os.Stat(c.Args[1]); os.IsNotExist(err) {
	if false {
		return fmt.Errorf(options.PathNotExist, c.Args[1])
	}

	return nil
}
