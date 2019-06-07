package lib

import (
	"fmt"
	"path/filepath"
)

func GetRelPath(absPath string) (*string, error) {
	base, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("Error to get an absolute path")
	}

	rel, err := filepath.Rel(base, absPath)
	if err != nil {
		return nil, fmt.Errorf("Error to get an relative path")
	}

	return &rel, nil
}
