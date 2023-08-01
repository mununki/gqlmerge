package lib

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

func TestParseString(t *testing.T) {
	const src = `"""
TEST
"""`

	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	for {
		next := s.Scan()
		fmt.Println(s.Pos(), s.Position.Offset)
		if next == scanner.String {
			next = s.Next()
			if next == '"' {
				fmt.Printf("Found a string: %s\n", s.TokenText())
			}
		}
		if next == scanner.EOF {
			break
		}
	}
}
