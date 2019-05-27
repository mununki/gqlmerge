package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Command for gqlmerge
type Command struct {
	Args   []string
	Paths  []string
	Output string
	Indent string
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
		Help:             Usage(),
		PathNotExist:     "❌ Path '%s' does not Exist",
		NotEnoughArgs:    "❌ Not enough arguments",
		OutputFileNeeded: "❌ Output file argument is needed",
		WrongOption:      "❌ Wrong options",
		Version:          "v0.2.1",
	}

	err := c.parseFlags()
	if err != nil {
		return fmt.Errorf("parse flags: %v", err)
	}

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

	c.Paths = c.Args[:argsCount-1]
	c.Output = c.Args[argsCount-1]

	// check passed paths is existing.
	// iter from 1 to argsCount-1 because the last
	// argument is an output file
	for _, path := range c.Paths {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf(options.PathNotExist, path)
		}
	}

	return nil
}

func (c *Command) parseFlags() (err error) {
	indent := flag.String("indent", "4s", flagIndentMsg)

	flag.Parse()

	// indent is never nil, so
	// dereference without checking
	c.Indent, err = convIndent(*indent)
	if err != nil {
		return
	}

	// flag.Parse() remove program's name (aka os.Args[0])
	// and parsed flags from os.Args, so
	// flag.Args() = os.Args - os.Args[0] - parsed flags
	// i.e flag.Args() contains only the paths to parse
	// and the output filename.
	c.Args = flag.Args()

	return
}

func convIndent(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("indent should be not empty")
	}

	nStr := s[:len(s)-1]

	// if s is correct, then all characters
	// except the last should be a number
	n, err := strconv.Atoi(nStr)
	if err != nil {
		// if nStr == "", then "n" should be 1, but
		// Atoi will error, so return the error only
		// if nStr != ""
		if nStr != "" {
			return "", err
		}

		n = 1
	}

	// get the last symbol which contains
	// type of indent: tab or space
	i := s[len(s)-1:]

	switch i {
	case "s":
		i = " "
	case "t":
		i = "\t"
	default:
		return "", fmt.Errorf(`unknown indent "%s"`, i)
	}

	return strings.Repeat(i, n), nil
}
