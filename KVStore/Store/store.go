package store

import (
	"time"
	"sync"
)

type KeyValueStore interface {
	Set (key string, value any, ttl time.Duration) 
	Get (key string) (any, error)
	Delete (key string) error
	Close()
}

type MemoryStore struct {
	mu sync.RWMutex
	items map[string]item
	stopWorker chan struct{}

}

func New(cleanupInterval time.Duration) *MemoryStore {
	s := &MemoryStore{
		items: make(map[string]item), 
		stopWorker: make(chan struct{}),
	}
	if cleanupInterval > 0 {
		go s.startCleanupWorker(cleanupInterval)
	}
	return s
}

func (s *MemoryStore) Set(key string, value any, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = item {
		value : value,
		expiration: exp,

	}
	
}

func (s *MemoryStore) Get(key string) (any,error) {
	s.mu.RLock()
	it,exists := s.items[key]
	s.mu.RUnlock()
	if !exists {
		return nil, ErrKeyNotFound
	}
	if it.isExpired() {
		s.Delete(key)
		return nil, ErrKeyExpired
	}
	return it.value,nil
}
func (s * MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _,exists := s.items[key]; !exists {
		return ErrKeyNotFound
	}
	delete(s.items,key)
	return nil
}
func (s * MemoryStore) startCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			s.deleteEXpired()
		case <- s.stopWorker:
			return
		}
	}

}

func (s *MemoryStore) deleteEXpired() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, it := range s.items {
		if it.expiration > 0 && now > it.expiration {
			delete(s.items,key)
		}
	}

	
}

func (s *MemoryStore) Close() {
	close(s.stopWorker)
}



