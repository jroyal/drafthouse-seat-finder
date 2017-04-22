package drafthouse

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once they hit the expiration
type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*CacheItem
}

// CacheItem represents a value in the cache map
type CacheItem struct {
	sync.RWMutex
	data    interface{}
	expires *time.Time
}

func (item *CacheItem) expired() bool {
	var value bool
	item.RLock()
	defer item.RUnlock()
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	return value
}

func newCacheItem(data interface{}, duration time.Duration) *CacheItem {
	expires := time.Now().Add(duration)
	return &CacheItem{
		data:    data,
		expires: &expires,
	}
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, data interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item := newCacheItem(data, cache.ttl)
	cache.items[key] = item
}

// Get is a thread-safe way to lookup items
func (cache *Cache) Get(key string) (data interface{}, found bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		data = item.data
		found = true
	}
	return
}

func (cache *Cache) cleanup() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
}

func (cache *Cache) startCleanupTimer() {
	ticker := time.NewTicker(cache.ttl)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.cleanup()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// NewCache is a helper to create instance of the Cache struct
func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*CacheItem{},
	}
	cache.startCleanupTimer()
	return cache
}
