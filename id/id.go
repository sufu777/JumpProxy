package id

import (
	"sync"
)

var lk sync.Mutex
var counter int64 = 0

func NextId() int64 {
	var inner int64
	lk.Lock()
	defer lk.Unlock()
	inner = counter
	inner = (inner + 1) % 1000000
	return inner
}
