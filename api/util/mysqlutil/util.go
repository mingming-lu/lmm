package mysqlutil

import "regexp"

var (
	patternDuplicateKey = regexp.MustCompile(`Error 1062: Duplicate entry '([-\w]+)' for key '(\w+)'`)
)

// CheckDuplicateKeyError reports if given error describes a MySQL Error 1062
func CheckDuplicateKeyError(err error) (key string, entry string, ok bool) {
	if err == nil {
		return key, entry, ok
	}
	matched := patternDuplicateKey.FindStringSubmatch(err.Error())
	if len(matched) == 3 {
		key = matched[2]
		entry = matched[1]
		ok = true
	}
	return key, entry, ok
}
