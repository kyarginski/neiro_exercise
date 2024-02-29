package repository

import (
	"log/slog"
	"sync"
	"time"
)

// KeyValueStore struct holds the data and provides synchronization.
type KeyValueStore struct {
	log   *slog.Logger
	store map[string]Value

	lock sync.RWMutex
}

// Value struct to hold the value and the expiration time.
type Value struct {
	value      string
	expiration time.Time
}

// NewKeyValueStore initializes a new store.
func NewKeyValueStore(log *slog.Logger) *KeyValueStore {
	kv := &KeyValueStore{
		log:   log,
		store: make(map[string]Value),
	}

	return kv
}

// Set method to add a new key-value pair with an optional TTL (in seconds).
func (kv *KeyValueStore) Set(key, value string, ttl int) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	expiration := time.Now().Add(time.Duration(ttl) * time.Second)
	kv.store[key] = Value{value: value, expiration: expiration}
}

// Get method to retrieve a value by key.
func (kv *KeyValueStore) Get(key string) (string, bool) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	val, exists := kv.store[key]
	if !exists || val.expiration.Before(time.Now()) {
		return "", false
	}
	return val.value, true
}

// Delete method to remove a key-value pair.
func (kv *KeyValueStore) Delete(key string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	delete(kv.store, key)
}

// ClearExpired runs in a separate goroutine to remove expired keys.
func (kv *KeyValueStore) ClearExpired() {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	for key, val := range kv.store {
		if val.expiration.Before(time.Now()) {
			kv.log.Debug("Delete expired key", "key", key)
			delete(kv.store, key)
		}
	}
}
