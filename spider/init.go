package spider

import "sync"

var m *sync.RWMutex

func init() {
	m = new(sync.RWMutex)
}
