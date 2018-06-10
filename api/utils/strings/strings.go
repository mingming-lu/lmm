package strings

import (
	"strconv"
)

func StrToInt(s string) (int, error) {
	n, err := StrToUint64(s)
	return int(n), err
}

func StrToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
