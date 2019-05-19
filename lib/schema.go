package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetSchema is to parse ./schema/**/*.graphql
func (sc *Schema) GetSchema(path string) {
	// FIX: is there any way to use a relative path?
	// currently, it works only with absolute path
	// in case of using a relative path such as '../schema', it spits out an error
	// the error says invalid memory or nil pointer deference.
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p == "" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(p) != ".graphql" {
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			fmt.Printf("[Error] There is an error to open %s", p)
			return nil
		}

		// TODO: split and get a only filename and print it to user
		// needs to handle in case of OS (windows / unix compatibles)
		sc.Files = append(sc.Files, file)

		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("🎉 Total %d *.graphql files found!\n", len(sc.Files))

	return
}
