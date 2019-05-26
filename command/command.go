package command

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command for gqlmerge
type Command struct {
	Args   []string
	Paths  []string
	Output string
	Ident  string
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
		Version:          "v0.2.1",
	}

	c.parseFlags()

	argsCount := len(c.Args)

	// check the number of args
	if argsCount == 0 {
		// no arg -> print help msg
		return fmt.Errorf(options.Help)
	}

	if argsCount == 1 {
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

	c.Paths = c.Args
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

func (c *Command) parseFlags() {
	ident := flag.String("ident", "4s", `
Flag "ident" defines the padding in the generated GraphQL scheme.

It follows the next pattern: ident={n}{i},
	* n - amount of idents
	* i - ident ("t" for tabs and "s" for spaces)
If "n" is not stated 1 will be used, 
so "--ident=1t" is equal to "--ident=t"
`)

	flag.Parse()

	// ident is never nil, so
	// dereference without checking
	c.Ident = *ident

	// flag.Parse() remove program's name (aka os.Args[0])
	// and parsed flags from os.Args, so
	// flag.Args() = os.Args - os.Args[0] - parsed flags
	// i.e flag.Args() contains only the paths to parse
	// and the output filename.
	c.Args = flag.Args()
}
