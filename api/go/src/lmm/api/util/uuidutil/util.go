package uuidutil

import (
	"lmm/api/util/stringutil"
	"regexp"

	"github.com/google/uuid"
)

var (
	patternRawHex = regexp.MustCompile(`(.{8})(.{4})(.{4})(.{4})(.{12})`)
	stdUUIDFormat = `$1-$2-$3-$4-$5`
)

// ParseString decodes s into a UUID or returns an error. The following formats are allowed
// with hyphen:    xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// without hyphen: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func ParseString(s string) (uuid.UUID, error) {
	if len(s) == 32 {
		s = patternRawHex.ReplaceAllString(s, stdUUIDFormat)
	}
	return uuid.Parse(s)
}

// New creates a new uuid string without '-',
// Ex. 7fc07047356443e991549c71332e7dfd
func New() string {
	return stringutil.ReplaceAll(uuid.New().String(), "-", "")
}
