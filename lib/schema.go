package lib

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2"
)

// GetSchema is to parse ./schema/**/*.graphql
func GetSchema(path string) string {
	rel, err := filepath.Rel("lib", path)
	if err != nil {
		panic(err)
	}

	box := packr.New("schema", rel)
	var schema strings.Builder

	box.Walk(func(p string, f packr.File) error {
		if p == "" {
			return nil
		}

		var err error
		if finfo, err := f.FileInfo(); err != nil {
			return err
		} else {
			if finfo.IsDir() {
				return nil
			}
		}

		if !strings.Contains(p, ".graphql") {
			return nil
		}

		n := p
		if strings.HasPrefix(p, "\\") || strings.HasPrefix(p, "/") {
			n = n[1:]
		}

		n = strings.Replace(n, "\\", "/", -1)

		s, err := box.FindString(p)
		if err != nil {
			return nil
		}

		schema.WriteString(s)

		return nil
	})

	return schema.String()
}
