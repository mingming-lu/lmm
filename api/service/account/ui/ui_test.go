package ui

import (
	"lmm/api/context/account/infra"
	"lmm/api/testing"
	"os"
)

var ui *UI

func TestMain(m *testing.M) {
	ui = New(infra.NewUserStorage(testing.DB()))

	code := m.Run()
	os.Exit(code)
}
