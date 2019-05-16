package main

import (
	"testing"
)

func TestUnique(t *testing.T) {
	var test = []int{1, 2, 3, 1, 2, 3, 4, 6, 2}

	t.Log(UniqueMap(test, t))
}

func UniqueMap(s []int, t *testing.T) []int {
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	t.Log(s)
	return s[:j]
}
