package sliceutil

import (
	"strings"

	"lmm/api/testing"
)

func TestAnyString(tt *testing.T) {
	t := testing.NewTester(tt)

	t.False(AnyString([]string{}, func(s string) bool {
		return s == "empty"
	}))

	t.False(AnyString([]string{"a", "b", "ab"}, func(s string) bool {
		return s == "not found"
	}))

	t.True(AnyString([]string{"Golang", "JavaScript", "Python"}, func(s string) bool {
		return s == "Python"
	}))

	t.True(AnyString([]string{"apple", "boy", "cat"}, func(s string) bool {
		return strings.Contains(`boy ♂ next ♂ door`, s)
	}))
}

func TestContainsString(tt *testing.T) {
	t := testing.NewTester(tt)

	t.False(ContainsString("empty", []string{}))
	t.False(ContainsString("Python", []string{"Golang", "Python3"}))
	t.True(ContainsString("Golang", []string{"Golang"}))
	t.True(ContainsString("Python", []string{"Golang", "Python"}))
}
