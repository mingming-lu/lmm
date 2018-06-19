package uuid

import (
	"lmm/api/utils/strings"

	"github.com/google/uuid"
)

func New() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
