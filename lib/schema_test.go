package lib

import (
	"path/filepath"
	"testing"
)

func TestGetSchema(t *testing.T) {
	rel, err := filepath.Rel("lib", "./schema")
	if err != nil {
		t.Error(err)
	}

	t.Log(rel)
}
