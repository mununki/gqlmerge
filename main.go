package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gql "github.com/mattdamon108/gqlmerge/lib"
)

func main() {
	cmd := gql.Command{Args: os.Args}
	if err := cmd.Check(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sc := gql.Schema{}

	// TODO : needs to improve to work with a relative path.
	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	s := sc.GetSchema(abs)

	l := gql.NewLexer(s)

	if len(sc.Files) > 0 {
		sc.ParseSchema(l)
		sc.UniqueMutation()
		sc.UniqueQuery()
		sc.UniqueTypeName()
		sc.UniqueScalar()
		sc.UniqueEnum()
		sc.UniqueInterface()
		sc.UniqueUnion()
		sc.UniqueInput()

		ms := gql.MergedSchema{}
		ss := ms.StitchSchema(&sc)

		bs := []byte(ss)
		err := ioutil.WriteFile(os.Args[2], bs, 0644)
		if err != nil {
			fmt.Printf("😱 Error in writing '%s' file", os.Args[2])
		}

		fmt.Printf("👍 Successfully generated '%s'", os.Args[2])
	} else {
		fmt.Printf("😳 Not found any *.graphql files in %s", os.Args[1])
	}
}
