package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSchema(t *testing.T) {
	err := filepath.Walk("../schema", func(p string, info os.FileInfo, err error) error {
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

		n := p
		if strings.HasPrefix(p, "\\") || strings.HasPrefix(p, "/") {
			n = n[1:]
		}

		n = strings.Replace(n, "\\", "/", -1)

		s, err := ioutil.ReadFile(p)
		if err != nil {
			fmt.Printf("[Error] There is an error to read %s", p)
			return nil
		}

		fmt.Printf("Found %s\n", string(s))

		return nil
	})
	if err != nil {
		panic(err)
	}

}
