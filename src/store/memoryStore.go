package store

import (
	"strings"
	"sync"
	"time"
)

type Item struct {
	Value     interface{}
	Timestamp time.Time
}

type memoryStore struct {
	items map[string]Item
	mutex sync.RWMutex
}

var (
	memoryStoreInstance *memoryStore
	memoryStoreOnce     sync.Once
)

func NewMemoryStore() MemoryStore {
	memoryStoreOnce.Do(func() {
		memoryStoreInstance = &memoryStore{
			items: make(map[string]Item),
		}
	})
	return memoryStoreInstance
}

// constructKey creates a full key based on actor and key.
func (s *memoryStore) constructKey(actor, key string) string {
	return actor + ":" + key
}

// Set stores a key-value pair for a specific actor.
func (s *memoryStore) Set(actor, key string, value interface{}) {
	fullKey := s.constructKey(actor, key)
	s.mutex.Lock()
	s.items[fullKey] = Item{Value: value, Timestamp: time.Now()}
	s.mutex.Unlock()
}

// Get retrieves an item by key for a specific actor.
func (s *memoryStore) Get(actor, key string) (interface{}, bool) {
	fullKey := s.constructKey(actor, key)
	s.mutex.RLock()
	item, exists := s.items[fullKey]
	s.mutex.RUnlock()
	return item.Value, exists
}

// GetAllForActor returns all items for a specific actor.
func (s *memoryStore) GetAllForActor(actor string) map[string]interface{} {
	actorPrefix := s.constructKey(actor, "")
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	actorData := make(map[string]interface{})
	for key, item := range s.items {
		if strings.HasPrefix(key, actorPrefix) {
			actorData[strings.TrimPrefix(key, actorPrefix)] = item.Value
		}
	}
	return actorData
}
