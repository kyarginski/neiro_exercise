package repository

import (
	"testing"

	"neiro/internal/lib/logger/sl"
)

func TestSetAndGet(t *testing.T) {
	logger := sl.SetupLogger("nop")
	kv := NewKeyValueStore(logger)
	key, value := "testKey", "testValue"
	ttl := 10 // Time to live in seconds

	kv.Set(key, value, ttl)
	got, exists := kv.Get(key)

	if !exists || got != value {
		t.Errorf("Get(%q) = %q, %t; want %q, true", key, got, exists, value)
	}
}

func TestGetExpired(t *testing.T) {
	logger := sl.SetupLogger("nop")
	kv := NewKeyValueStore(logger)
	key, value := "expiredKey", "testValue"
	ttl := -1 // Set an expired TTL

	kv.Set(key, value, ttl)
	_, exists := kv.Get(key)

	if exists {
		t.Errorf("Get(%q) = %t; want false", key, exists)
	}
}

func TestDelete(t *testing.T) {
	logger := sl.SetupLogger("nop")
	kv := NewKeyValueStore(logger)
	key, value := "deleteKey", "testValue"
	kv.Set(key, value, 10)

	kv.Delete(key)
	_, exists := kv.Get(key)

	if exists {
		t.Errorf("Delete(%q) succeeded, but Get(%q) = true; want false", key, key)
	}
}

func TestClearExpired(t *testing.T) {
	logger := sl.SetupLogger("nop")
	kv := NewKeyValueStore(logger)
	expiredKey, validKey := "expiredKey", "validKey"
	kv.Set(expiredKey, "expiredValue", -1) // Expired
	kv.Set(validKey, "validValue", 10)     // Valid

	kv.ClearExpired()

	_, expiredExists := kv.Get(expiredKey)
	_, validExists := kv.Get(validKey)

	if expiredExists {
		t.Errorf("ClearExpired did not remove expired key %q", expiredKey)
	}

	if !validExists {
		t.Errorf("ClearExpired incorrectly removed valid key %q", validKey)
	}
}
