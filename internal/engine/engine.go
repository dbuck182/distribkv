package engine

import "sync"

type Engine struct {
	table map[string][]byte
	mu    sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		table: make(map[string][]byte),
	}
}

// Locks the table and returns the key when it is its turn
// When it has the lock

func (e *Engine) Get(key string) []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.table[key]
}

// Sets a new key

func (e *Engine) Set(key string, data []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.table[key] = data
}

func (e *Engine) Del(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.table, key)
}
