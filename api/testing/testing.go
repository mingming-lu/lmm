package testing

import (
	"lmm/api/storage"
	"math/rand"
	"sync"
	"time"
)

var (
	cache *storage.Cache
	db    *storage.DB
)

func init() {
	db = storage.NewDB()
	cache = storage.NewCacheEngine()

	rand.Seed(time.Now().UnixNano())
	InitTableAll()
}

func DB() *storage.DB {
	return db
}

func Cache() *storage.Cache {
	return cache
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
