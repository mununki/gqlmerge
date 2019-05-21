package command

import (
	"fmt"
	"os"
	"strings"
)

// Command for gqlmerge
type Command struct {
	Args   []string
	Paths  []string
	Output string
}

type Options struct {
	Help             string
	PathNotExist     string
	NotEnoughArgs    string
	OutputFileNeeded string
	WrongOption      string
	Version          string
}

func (c *Command) Check() error {
	options := Options{
		Help: `👋 'gqlmerge' is the tool to merge & stitch GraphQL files and generate a GraphQL schema
Author : Woonki Moon <woonki.moon@gmail.com>

Usage:	gqlmerge [PATH] [OUTPUT]

e.g.

	gqlmerge ./schema schema.graphql

Options:

	-v	: check the version
	-h	: help
`,
		PathNotExist:     "❌ Path '%s' does not Exist",
		NotEnoughArgs:    "❌ Not enough arguments",
		OutputFileNeeded: "❌ Output file argument is needed",
		WrongOption:      "❌ Wrong options",
		Version:          "v0.2.0",
	}

	argsCount := len(c.Args)

	// check the number of args
	if argsCount <= 1 {
		// no arg -> print help msg
		return fmt.Errorf(options.Help)
	}

	if argsCount == 2 {
		if strings.HasPrefix(c.Args[1], "-") {
			switch c.Args[1] {
			case "-v":
				return fmt.Errorf(options.Version)
			case "-h":
				return fmt.Errorf(options.Help)
			default:
				return fmt.Errorf(options.WrongOption)
			}
		}

		return fmt.Errorf(options.OutputFileNeeded)
	}

	c.Paths = c.Args[1 : argsCount-1]
	c.Output = c.Args[argsCount-1]

	// check passed paths is existing.
	// iter from 1 to argsCount-1 because the last
	// argument is an output file
	for _, path := range c.Paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf(options.PathNotExist, path)
		}
	}

	return nil
}
