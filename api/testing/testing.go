package testing

import (
	"lmm/api/storage"
	"math/rand"
	"sync"
	"time"
)

var (
	cacheEngine *storage.Cache
	dbEngine    *storage.DB
)

func init() {
	dbEngine = storage.NewDB()
	cacheEngine = storage.NewCacheEngine()

	rand.Seed(time.Now().UnixNano())
	InitTableAll()
}

func DB() *storage.DB {
	return dbEngine
}

func Cache() *storage.Cache {
	return cacheEngine
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
