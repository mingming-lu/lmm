package testing

import (
	"testing"

	"github.com/pkg/errors"
)

func TestIs(t *testing.T) {
	tester := NewTester(t)
	tester.Is(1, 1)
	tester.Is(true, true)
	tester.Is("abc", "abc")
}

func TestNot(t *testing.T) {
	tester := NewTester(t)
	tester.Not(1, 2)
	tester.Not(true, false)
	tester.Not(1, "1")
	tester.Not(tester, t)
}

func TestError(t *testing.T) {
	tester := NewTester(t)
	tester.Error(errors.New(""))
	tester.Error(errors.New("msg"))
}

func TestNoError(t *testing.T) {
	tester := NewTester(t)
	tester.NoError(nil)
}
