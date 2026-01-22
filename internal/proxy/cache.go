package proxy

import (
	"container/list"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	FilePath   string
	Size       int64
	LastAccess time.Time
}

type LRUCache struct {
	maxSize     int64
	currentSize int64
	items       map[string]*list.Element
	list        *list.List
	mutex       sync.RWMutex
	cacheDir    string
}

var imageCache *LRUCache

func InitCache(maxSize int64, cacheDir string) error {
	var err error
	imageCache, err = NewLRUCache(maxSize, cacheDir)
	return err
}

func NewLRUCache(maxSize int64, cacheDir string) (*LRUCache, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, err
	}
	return &LRUCache{
		maxSize:  maxSize,
		items:    make(map[string]*list.Element),
		list:     list.New(),
		cacheDir: cacheDir,
	}, nil
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		item := elem.Value.(*CacheItem)
		item.LastAccess = time.Now()
		data, err := os.ReadFile(item.FilePath)
		if err != nil {
			return nil, false
		}
		return data, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	size := int64(len(data))
	for c.currentSize+size > c.maxSize && c.list.Len() > 0 {
		c.evictOldest()
	}

	if size > c.maxSize {
		return fmt.Errorf("file too large")
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
	filePath := filepath.Join(c.cacheDir, hash+".jpg")
	if err := os.WriteFile(filePath, data, 0o600); err != nil {
		return err
	}

	item := &CacheItem{
		Key:        key,
		FilePath:   filePath,
		Size:       size,
		LastAccess: time.Now(),
	}
	elem := c.list.PushFront(item)
	c.items[key] = elem
	c.currentSize += size
	return nil
}

func (c *LRUCache) evictOldest() {
	elem := c.list.Back()
	if elem == nil {
		return
	}
	item := elem.Value.(*CacheItem)
	os.Remove(item.FilePath)
	c.list.Remove(elem)
	delete(c.items, item.Key)
	c.currentSize -= item.Size
}
