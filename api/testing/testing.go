package testing

import (
	"sync"
)

func init() {
	InitTableAll()
}

var mutex sync.Mutex

func Lock() {
	mutex.Lock()
}

func Unlock() {
	mutex.Unlock()
}
