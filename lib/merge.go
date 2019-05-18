package lib

import (
	"sync"
)

func Merge(path string) *string {

	sc := Schema{}
	// at this moment, path should be an absolute path
	s := sc.GetSchema(path)

	l := NewLexer(s)

	if len(sc.Files) > 0 {
		sc.ParseSchema(l)

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
	return nil
}
