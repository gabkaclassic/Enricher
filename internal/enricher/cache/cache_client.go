package cache

type CacheClient interface {

	Get(key string) ([]byte, error)
	SetWithTTL(key string, value []byte, ttl int) (error)
	Set(key string, value []byte) (error)
	Delete(keys...string) (int64, error)
	Clean() error
}