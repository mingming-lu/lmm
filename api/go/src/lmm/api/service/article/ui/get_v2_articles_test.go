package ui

import (
	"lmm/api/testing"
)

func TestListArticleV2(tt *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	// TODO GET /v2/articles
}
