package testing

import (
	"lmm/api/storage"
	"math/rand"
	"sync"
	"time"
)

var db *storage.DB

func init() {
	db = storage.NewDB()
	rand.Seed(time.Now().UnixNano())
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
