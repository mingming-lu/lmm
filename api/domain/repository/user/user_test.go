package user

import (
	"testing"

	model "lmm/api/domain/model/user"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	repo := New()
	m := model.New("foobar", "1234")
	_, err := repo.Save(m)
	assert.NoError(t, err)
}
