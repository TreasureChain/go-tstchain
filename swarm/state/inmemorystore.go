package state

import (
	"encoding"
	"encoding/json"
	"sync"
)

// InmemoryStore is the reference implementation of Store interface that is supposed
// to be used in tests.
type InmemoryStore struct {
	db map[string][]byte
	mu sync.RWMutex
}

// NewInmemoryStore returns a new instance of InmemoryStore.
func NewInmemoryStore() *InmemoryStore {
	return &InmemoryStore{
		db: make(map[string][]byte),
	}
}

// Get retrieves a value stored for a specific key. If there is no value found,
// ErrNotFound is returned.
func (s *InmemoryStore) Get(key string, i interface{}) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bytes, ok := s.db[key]
	if !ok {
		return ErrNotFound
	}

	unmarshaler, ok := i.(encoding.BinaryUnmarshaler)
	if !ok {
		return json.Unmarshal(bytes, i)
	}

	return unmarshaler.UnmarshalBinary(bytes)
}

// Put stores a value for a specific key.
func (s *InmemoryStore) Put(key string, i interface{}) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bytes := []byte{}

	marshaler, ok := i.(encoding.BinaryMarshaler)
	if !ok {
		if bytes, err = json.Marshal(i); err != nil {
			return err
		}
	} else {
		if bytes, err = marshaler.MarshalBinary(); err != nil {
			return err
		}
	}

	s.db[key] = bytes
	return nil
}

// Delete removes value stored under a specific key.
func (s *InmemoryStore) Delete(key string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[key]; !ok {
		return ErrNotFound
	}
	delete(s.db, key)
	return nil
}

// Close does not do anything.
func (s *InmemoryStore) Close() error {
	return nil
}
