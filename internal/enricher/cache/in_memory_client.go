package cache

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryCacheClient struct {
	memory sync.Map
}

type CacheItem struct {
	Value      []byte
	Expiration time.Time
}

func NewInMemoryCacheClient() *InMemoryCacheClient {

	return &InMemoryCacheClient{}
}

func (cache *InMemoryCacheClient) Get(key string) ([]byte, error) {

	value, exists := cache.memory.Load(key)

	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	entry := value.(CacheItem)
	expired := entry.Expiration

	if !expired.IsZero() && expired.Before(time.Now()) {
		cache.Delete(key)
		return nil, fmt.Errorf("key '%s' expired", key)
	}

	return entry.Value, nil
}

func (cache *InMemoryCacheClient) Set(key string, value []byte) error {

	cache.memory.Store(key, CacheItem{
		Value:      value,
		Expiration: time.Time{},
	})

	return nil
}

func (cache *InMemoryCacheClient) SetWithTTL(key string, value []byte, ttl int) error {

	cache.memory.Store(key, CacheItem{
		Value:      value,
		Expiration: time.Now().Add(time.Duration(ttl * int(time.Second))),
	})
	return nil
}

func (cache *InMemoryCacheClient) Delete(keys ...string) (int64, error) {

	var deleteCount int64
	for _, k := range keys {
		if _, exists := cache.memory.Load(k); exists {
			cache.memory.Delete(k)
			deleteCount++
		}
	}

	return deleteCount, nil
}

func (cache *InMemoryCacheClient) Clean() error {

	cache.memory.Clear()

	return nil
}
