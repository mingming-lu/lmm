package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticle(t *testing.T) {
	article := Article{}

	assert.False(t, article.Published())
}
