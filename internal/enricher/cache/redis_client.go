package cache

import (
	"github.com/redis/go-redis/v9"
	"time"
	"context"
)

type RedisCacheClient struct {
	client *redis.Client
}

func NewRedisCacheClient(address string, password string, db int) *RedisCacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: db,
	})

	return &RedisCacheClient{client: rdb}
}

func (cache *RedisCacheClient) Get(key string) ([]byte, error) {

	ctx := context.Background()

	stringValue, err := cache.client.Get(ctx, key).Result()

	return []byte(stringValue), err
}

func (cache *RedisCacheClient) Set(key string, value []byte) (error) {

	ctx := context.Background()

	_, err := cache.client.Set(ctx, key, value, 0).Result()
	return err
}

func (cache *RedisCacheClient) SetWithTTL(key string, value []byte, ttl int) (error) {

	ctx := context.Background()

	_, err := cache.client.Set(ctx, key, value, time.Duration(ttl)).Result()
	return err
}

func (cache *RedisCacheClient) Delete(keys...string) (int64, error) {

	ctx := context.Background()

	return cache.client.Del(ctx, keys...).Result()
}

func (cache *RedisCacheClient) Clean() (error) {

	ctx := context.Background()

	_, err := cache.client.FlushAllAsync(ctx).Result()

	return err
}