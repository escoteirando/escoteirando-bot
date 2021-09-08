package utils

import (
	"sync"
	"time"
)

type UniqueSequence struct {
	sync.RWMutex
	current int64
}

func NewUniqueSequence() UniqueSequence {
	return UniqueSequence{
		current: 0,
	}
}
func (us *UniqueSequence) GetNext() uint {
	us.Lock()
	n := time.Now().Unix()
	if us.current == n {
		n = us.current + 1
	}
	us.current = n
	us.Unlock()
	return uint(us.current)
}
