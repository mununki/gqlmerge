package lib

import (
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	var src = `
	type Query {
		checkIfExists(userId: ID!, name: String): CheckIfExistsResponse!
	}

	directive @goModel(
		model: String, models: [String!]
	) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION
	
	"""
	TEST
	"""
	interface Node @goModel(
		" description "
		model: "todo/ent.Noder", models: ["a", "b"]) {
		id: ID!
	}
	
	enum Color @goModel(model: "backend/ent/color.Color") {
		Blue @ignore(if: isError)
		Red
	}

	type User implements Node & Profile {
		id: ID!
		email: String!
		fullName: String!
	}

	type CheckIfExistsResponse {
		ok: Boolean!
		error: String
		user: [User]!
	}
`

	s := Schema{}
	p := NewParser(strings.NewReader(src), "")
	s.Parse(p)
}
