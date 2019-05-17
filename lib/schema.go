package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetSchema is to parse ./schema/**/*.graphql
func (sc *Schema) GetSchema(path string) string {
	var schema strings.Builder

	// FIX: is there any way to use a relative path?
	// currently, it works only with absolute path
	// in case of using a relative path such as '../schema', it spits out an error
	// the error says invalid memory or nil pointer deference.
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if p == "" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.Contains(p, ".graphql") {
			return nil
		}

		s, err := ioutil.ReadFile(p)
		if err != nil {
			fmt.Printf("[Error] There is an error to read %s", p)
			return nil
		}

		// TODO: split and get a only filename and print it to user
		// needs to handle in case of OS (windows / unix compatibles)
		sc.Files = append(sc.Files, p)

		schema.Write(s)

		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("🎉 Total %d *.graphql files found!\n", len(sc.Files))

	return schema.String()
}
