package lib

import (
	"strings"
	"sync"
	"testing"
)

func TestMerge(t *testing.T) {
	var src = `
	schema {
		query: Query
		mutation: Mutation
	}

	schema {
		query: Query
		mutation: Mutation2
		subscription: Subscription
	}

	type Query {
		checkIfExists(userId: ID!, name: String): CheckIfExistsResponse!
	}

	type Query {
		getMyProfile: UserResponse!
	}
`

	s := Schema{}
	p := NewParser(strings.NewReader(src), "")
	s.Parse(p)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go s.mergeSchemaDefinition(&wg)

	wg.Wait()
}
