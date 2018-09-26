package ui

import (
	"os"

	"lmm/api/service/account/infra"
	"lmm/api/testing"
)

var ui *UI

func TestMain(m *testing.M) {
	ui = New(infra.NewUserStorage(testing.DB()))

	code := m.Run()
	os.Exit(code)
}
