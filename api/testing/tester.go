package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (t *Tester) NotNil(o interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotNil(t, o, msgAndArgs...)
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

func (t *Tester) IsError(expected, actual error, msgAndArgs ...interface{}) bool {
	if !assert.Error(t, expected) {
		return assert.FailNow(t, fmt.Sprintf("expected error but got %+v", expected))
	}
	return assert.EqualError(t, actual, expected.Error(), msgAndArgs...)
}

func (t *Tester) IsErrorMsg(expected string, actual error, msgAndArgs ...interface{}) bool {
	return assert.EqualError(t, actual, expected, msgAndArgs...)
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

func (t *Tester) Regexp(expected, actual string, msgAndArgs ...interface{}) bool {
	return assert.Regexp(t, expected, actual, msgAndArgs...)
}

func (t *Tester) JSON(expected interface{}, actual string, msgAndArgs ...interface{}) bool {
	b, err := json.Marshal(expected)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	return assert.JSONEq(t, string(b), actual)
}

// Are compares two collections' values
func (t *Tester) Are(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.EqualValues(t, expected, actual, msgAndArgs...)
}
