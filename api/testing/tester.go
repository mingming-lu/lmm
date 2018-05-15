package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type T = testing.T

type Tester struct {
	*T
}

func NewTester(t *T) *Tester {
	return &Tester{T: t}
}

func (t *Tester) Nil(o interface{}, magAndArgs ...interface{}) bool {
	return assert.Nil(t, o, magAndArgs...)
}

// Is compare both their types and referenced value of two interfaces.
// Returns true only if they hold same type and same value,
// no matter their memory addresses are same or not.
func (t *Tester) Is(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Exactly(t, expected, actual, msgAndArgs...)
}

func (t *Tester) Not(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotEqual(t, expected, actual, msgAndArgs...)
}

func (t *Tester) Error(err error, msgAndArgs ...interface{}) bool {
	return assert.Error(t, err, msgAndArgs...)
}

func (t *Tester) NoError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(t, err, msgAndArgs...)
}

func (t *Tester) Isa(expectedType, o interface{}, msgAndArgs ...interface{}) bool {
	return assert.IsType(t, expectedType, o, msgAndArgs...)
}
