package testing

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type M = testing.M

type TestRunner struct {
	m            *M
	code         *int
	setupFunc    func()
	teardownFunc func()
}

func NewTestRunner(m *M) *TestRunner {
	code := 0
	return &TestRunner{
		m:            m,
		code:         &code,
		setupFunc:    func() {},
		teardownFunc: func() {},
	}
}

func (r *TestRunner) Setup(setupFunc func()) *TestRunner {
	r.setupFunc = setupFunc
	return r
}

func (r *TestRunner) Teardown(tearDownFunc func()) *TestRunner {
	r.teardownFunc = tearDownFunc
	return r
}

func (r *TestRunner) Run() {
	r.setupFunc()
	defer os.Exit(r.m.Run())
	defer r.teardownFunc()
}
