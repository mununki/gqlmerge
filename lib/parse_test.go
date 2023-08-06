package lib

import (
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	var src = `
	directive @goModel(
		model: String, models: [String!]
	) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION
	
	"""
	TEST
	"""
	interface Node @goModel(model: "todo/ent.Noder", models: ["a", "b"]) {
		id: ID! # TEST
	}
	
	enum Color @goModel(model: "backend/ent/color.Color") {
		Blue @ignore(if: isError)
		Red
	}	
`

	s := Schema{}
	p := NewParser(strings.NewReader(src), "")
	s.Parse(p)
}
