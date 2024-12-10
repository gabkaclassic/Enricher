package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCacheClient struct {
	client *redis.Client
}

func NewRedisCacheClient(address string, password string, db int) *RedisCacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisCacheClient{client: rdb}
}

func (cache *RedisCacheClient) Get(key string) (interface{}, error) {

	ctx := context.Background()

	value, err := cache.client.Get(ctx, key).Bytes()

	return value, err
}

func (cache *RedisCacheClient) Set(key string, value interface{}) error {

	ctx := context.Background()

	_, err := cache.client.Set(ctx, key, value, 0).Result()
	return err
}

func (cache *RedisCacheClient) SetWithTTL(key string, value interface{}, ttl int) error {

	ctx := context.Background()

	_, err := cache.client.Set(ctx, key, value, time.Duration(ttl*int(time.Second))).Result()
	return err
}

func (cache *RedisCacheClient) Delete(keys ...string) (int64, error) {

	ctx := context.Background()

	return cache.client.Del(ctx, keys...).Result()
}

func (cache *RedisCacheClient) Clean() error {

	ctx := context.Background()

	_, err := cache.client.FlushAllAsync(ctx).Result()

	return err
}
