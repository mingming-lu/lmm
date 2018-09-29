package testing

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type M = testing.M
type T = testing.T
type Tester struct {
	*T
}

// NewTester creates a new Tester
func NewTester(t *T) *Tester {
	return &Tester{T: t}
}

// Nil asserts a interface is nil
func (t *Tester) Nil(o interface{}, msgAndArgs ...interface{}) bool {
	return assert.Nil(t, o, msgAndArgs...)
}

// NotNil asserts a interface is not nil
func (t *Tester) NotNil(o interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotNil(t, o, msgAndArgs...)
}

// True asserts a boolean value is true
func (t *Tester) True(value bool, msgAndArgs ...interface{}) bool {
	return assert.True(t, value, msgAndArgs...)
}

// False asserts a boolean value is false
func (t *Tester) False(value bool, msgAndArgs ...interface{}) bool {
	return assert.False(t, value, msgAndArgs)
}

// Is compare both their types and referenced value of two interfaces.
// Returns true only if they hold same type and same value,
// no matter their memory addresses are same or not.
func (t *Tester) Is(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Exactly(t, expected, actual, msgAndArgs...)
}

// Not asserts expected is not equal to given actual
func (t *Tester) Not(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.NotEqual(t, expected, actual, msgAndArgs...)
}

// Error asserts the given error is not nil
func (t *Tester) Error(err error, msgAndArgs ...interface{}) bool {
	return assert.Error(t, err, msgAndArgs...)
}

// NoError asserts the given error is nil
func (t *Tester) NoError(err error, msgAndArgs ...interface{}) bool {
	return assert.NoError(t, err, msgAndArgs...)
}

// IsError asserts actual error is equal to expected error, by comparing Error()
func (t *Tester) IsError(expected, actual error, msgAndArgs ...interface{}) bool {
	if !assert.Error(t, expected) {
		return assert.FailNow(t, fmt.Sprintf("expected error but got %+v", expected))
	}
	return assert.EqualError(t, actual, expected.Error(), msgAndArgs...)
}

// Isa asserts two interface share the same type
func (t *Tester) Isa(expectedType, o interface{}, msgAndArgs ...interface{}) bool {
	return assert.IsType(t, expectedType, o, msgAndArgs...)
}

// Regexp asserts s matches regexp
func (t *Tester) Regexp(regexp, s string, msgAndArgs ...interface{}) bool {
	return assert.Regexp(t, regexp, s, msgAndArgs...)
}

// JSON asserts a interface can be encoded to the given string
func (t *Tester) JSON(expected interface{}, actual string, msgAndArgs ...interface{}) bool {
	b, err := json.Marshal(expected)
	if err != nil {
		return assert.Fail(t, err.Error())
	}
	return assert.JSONEq(t, string(b), actual)
}

// Are compares two collections' values
func (t *Tester) Are(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.EqualValues(t, expected, actual, msgAndArgs...)
}

// NotPanic asserts given function would not panics
func (t *Tester) NotPanic(f func(), msgAndArgs ...interface{}) {
	assert.NotPanics(t, f, msgAndArgs...)
}
