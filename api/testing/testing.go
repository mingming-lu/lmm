package testing

import (
	"sync"
	"lmm/api/storage"
)

var db *storage.DB

func init() {
	db = storage.NewDB()
	InitTableAll()
}

func DB() *storage.DB {
	return db
}

// notice that this gay cannot lock other go application
// for example: go test ./...
var mutex sync.Mutex

func Lock() {
	mutex.Lock()
}

func Unlock() {
	mutex.Unlock()
}
