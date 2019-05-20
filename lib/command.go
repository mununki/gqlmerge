package lib

import (
	"fmt"
	"os"
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
	WrongOption      string
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
		WrongOption:      "❌ Wrong options",
		Version:          "v0.1.4",
	}

	// check the number of args
	if len(c.Args) <= 1 {
		// no arg -> print help msg
		return fmt.Errorf(options.Help)
	}

	if len(c.Args) == 2 {
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

	// check first arg, path is existing
	if _, err := os.Stat(c.Args[1]); os.IsNotExist(err) {
		return fmt.Errorf(options.PathNotExist, c.Args[1])
	}

	return nil
}
