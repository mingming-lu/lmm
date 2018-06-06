package repository

import (
	"lmm/api/db"
	"regexp"
)

var (
	ErrPatternDuplicate = regexp.MustCompile(`Error 1062: Duplicate entry '([-\w]+)' for key '(\w+)'`)
)

type Repository interface {
	DB() *db.DB
}

type Default struct {
}

func (repo *Default) DB() *db.DB {
	return db.New()
}

func CheckErrorDuplicate(errMsg string) (key string, entry string, ok bool) {
	matched := ErrPatternDuplicate.FindStringSubmatch(errMsg)
	if len(matched) == 3 {
		key = matched[2]
		entry = matched[1]
		ok = true
	}
	return key, entry, ok
}
