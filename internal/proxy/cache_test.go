package proxy

import (
	"testing"
)

func TestLRUCache_SetGetEvict(t *testing.T) {
	cacheDir := t.TempDir()
	cache, err := NewLRUCache(100, cacheDir)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	data1 := []byte("12345")
	if err := cache.Set("key1", data1); err != nil {
		t.Fatalf("Failed to set key1: %v", err)
	}

	got, ok := cache.Get("key1")
	if !ok {
		t.Fatalf("Expected key1 to exist")
	}
	if string(got) != string(data1) {
		t.Errorf("Expected %s, got %s", string(data1), string(got))
	}

	data2 := make([]byte, 100)
	for i := range data2 {
		data2[i] = 'a'
	}
	if err := cache.Set("key2", data2); err != nil {
		t.Fatalf("Failed to set key2: %v", err)
	}

	if _, ok := cache.Get("key1"); ok {
		t.Errorf("Expected key1 to be evicted")
	}

	if _, ok := cache.Get("key2"); !ok {
		t.Errorf("Expected key2 to exist")
	}
}
