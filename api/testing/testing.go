package testing

import (
	"sync"
)

func init() {
	InitTableAll()
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
