package testing

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type M = testing.M
type T = testing.T

type Tester struct {
	*T
}

func NewTester(t *T) *Tester {
	return &Tester{T: t}
}

func (t *Tester) Nil(o interface{}, msgAndArgs ...interface{}) bool {
	return assert.Nil(t, o, msgAndArgs...)
}

func (t *Tester) True(value bool, msgAndArgs ...interface{}) bool {
	return assert.True(t, value, msgAndArgs...)
}

func (t *Tester) False(value bool, msgAndArgs ...interface{}) bool {
	return assert.False(t, value, msgAndArgs)
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

// Output matches output of given function with regexp
func (t *Tester) Output(expectedRegexp string, f func(), msgAndArgs ...interface{}) bool {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)

	return assert.Regexp(t, expectedRegexp, buf.String(), msgAndArgs...)
}
