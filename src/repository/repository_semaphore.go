package repository

import "sync"


type repositorySemaphore struct {
	mu sync.RWMutex
	CanWrite bool
}

var sema = repositorySemaphore{
	CanWrite: true,
}

func CanIWrite() {
	sema.mu.Lock()
}

func YouCanWrite(){
	sema.mu.Unlock()
}


