package main

import (
	"fmt"
	"os"

	"github.com/mununki/gqlmerge/command"
	gql "github.com/mununki/gqlmerge/lib"
)

func main() {
	cmd := command.Command{Args: os.Args}
	if err := cmd.Check(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// TODO : needs to improve to work with a relative path.

	ss := gql.Merge(cmd.Indent, cmd.Paths...)

	if ss != nil {
		bs := []byte(*ss)
		err := os.WriteFile(cmd.Output, bs, 0644)
		if err != nil {
			fmt.Printf("😱 Error in writing '%s' file", cmd.Output)
			return
		}

		fmt.Printf("👍 Successfully generated '%s'\n", cmd.Output)
	} else {
		fmt.Printf("😳 Not found any GraphQL files in %v\n", cmd.Paths)
	}
}
