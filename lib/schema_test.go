package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSchema(t *testing.T) {
	err := filepath.Walk("../test", func(p string, info os.FileInfo, err error) error {
		if p == "" {
			return nil
		}

		if info.IsDir() {
			fmt.Println(p)
			return nil
		}

		if !strings.Contains(p, ".graphql") {
			return nil
		}

		s, e := os.ReadFile(p)
		if e != nil {
			fmt.Printf("[Error] There is an error to read %s", p)
			return nil
		}

		fmt.Printf("Found\n%s\n", string(s))

		return nil
	})
	if err != nil {
		panic(err)
	}

}

func TestMergeDirectives(t *testing.T) {
	str := make([]string, 0)
	a := []*Directive{
		{
			Name:          "talkable",
			DirectiveArgs: []*DirectiveArg{},
			Descriptions:  &str,
		},
	}
	b := []*Directive{
		{
			Name:          "talkable",
			DirectiveArgs: []*DirectiveArg{},
			Descriptions:  &str,
		},
	}
	c := []*Directive{
		{
			Name:          "walkable",
			DirectiveArgs: []*DirectiveArg{},
			Descriptions:  &str,
		},
	}
	ds := mergeDirectives(a, b)

	if len(ds) == 0 {
		t.Fatal("should be more than 0")
	}

	ds = mergeDirectives(a, c)

	if len(ds) == 0 || len(ds) == 1 {
		t.Fatal("should be more than 1")
	}
}
