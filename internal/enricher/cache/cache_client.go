package cache

type CacheClient interface {
	Get(key string) (interface{}, error)
	SetWithTTL(key string, value interface{}, ttl int) error
	Set(key string, value interface{}) error
	Delete(keys ...string) (int64, error)
	Clean() error
}
