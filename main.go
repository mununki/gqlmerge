package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	cmd := Command{args: os.Args}
	if err := cmd.Check(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	s := GetSchema(os.Args[1])
	l := NewLexer(s)

	sc := Schema{}
	sc.ParseSchema(l)
	sc.UniqueMutation()
	sc.UniqueQuery()
	sc.UniqueTypeName()
	sc.UniqueScalar()
	sc.UniqueEnum()
	sc.UniqueInterface()
	sc.UniqueUnion()
	sc.UniqueInput()

	ms := MergedSchema{}
	ss := ms.StitchSchema(&sc)

	bs := []byte(ss)
	err := ioutil.WriteFile(os.Args[2], bs, 0644)
	if err != nil {
		fmt.Printf("Error in writing '%s' file", os.Args[2])
	}

	fmt.Printf("Successfully generated '%s'", os.Args[2])
}
