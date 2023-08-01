package lib

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

func TestParseString(t *testing.T) {
	const src = `(model: "todo/ent.Noder")`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == scanner.String {
			fmt.Printf("Found a string: %s\n", s.TokenText())
		}
	}
}
