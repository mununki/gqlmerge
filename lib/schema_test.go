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
