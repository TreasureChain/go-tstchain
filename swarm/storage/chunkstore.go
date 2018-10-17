package storage

import (
	"context"
	"sync"
)

/*
ChunkStore interface is implemented by :

- MemStore: a memory cache
- DbStore: local disk/db store
- LocalStore: a combination (sequence of) memStore and dbStore
- NetStore: cloud storage abstraction layer
- FakeChunkStore: dummy store which doesn't store anything just implements the interface
*/
type ChunkStore interface {
	Put(context.Context, *Chunk) // effectively there is no error even if there is an error
	Get(context.Context, Address) (*Chunk, error)
	Close()
}

// MapChunkStore is a very simple ChunkStore implementation to store chunks in a map in memory.
type MapChunkStore struct {
	chunks map[string]*Chunk
	mu     sync.RWMutex
}

func NewMapChunkStore() *MapChunkStore {
	return &MapChunkStore{
		chunks: make(map[string]*Chunk),
	}
}

func (m *MapChunkStore) Put(ctx context.Context, chunk *Chunk) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chunks[chunk.Addr.Hex()] = chunk
	chunk.markAsStored()
}

func (m *MapChunkStore) Get(ctx context.Context, addr Address) (*Chunk, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	chunk := m.chunks[addr.Hex()]
	if chunk == nil {
		return nil, ErrChunkNotFound
	}
	return chunk, nil
}

func (m *MapChunkStore) Close() {
}
