package lib

import (
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	var src = `
	schema {
		query: Query
		mutation: Mutation
	}

	type Query {
		checkIfExists(userId: ID!, name: String): CheckIfExistsResponse! # TEST
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
		id: ID! # TEST
	}
	
	enum Color @goModel(model: "backend/ent/color.Color") {
		Blue @ignore(if: isError) # TEST
		Red # TEST
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

	scalar DateTime # TEST

	union Response = Success | Failure
`

	s := Schema{}
	p := NewParser(strings.NewReader(src), "")
	s.Parse(p)
}
