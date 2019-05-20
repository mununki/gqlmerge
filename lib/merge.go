package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func Merge(path string) *string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	sc := Schema{}
	// at this moment, path should be an absolute path
	sc.GetSchema(abs)

	if len(sc.Files) == 0 {
		return nil
	}

	for _, file := range sc.Files {
		l := NewLexer(file)
		sc.ParseSchema(l)
	}

	var wg sync.WaitGroup

	wg.Add(8)

	go sc.UniqueMutation(&wg)
	go sc.UniqueQuery(&wg)
	go sc.UniqueTypeName(&wg)
	go sc.UniqueScalar(&wg)
	go sc.UniqueEnum(&wg)
	go sc.UniqueInterface(&wg)
	go sc.UniqueUnion(&wg)
	go sc.UniqueInput(&wg)

	wg.Wait()

	ms := MergedSchema{}
	ss := ms.StitchSchema(&sc)
	return &ss
}
