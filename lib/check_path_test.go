package lib

import (
	"os"
	"testing"
)

func TestCheckPath(t *testing.T) {
	if _, err := os.Stat("./schem"); os.IsNotExist(err) {
		t.Logf("not Exist")
	}
}
