package redis

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/go-redis/redis"
	"github.com/rtgnx/storage/kv"
)

type RedisStore struct {
	redis *redis.Client
}

func New(url string) (kv.KVStore, error) {
	x := &RedisStore{}
	opt, err := redis.ParseURL(url)

	if err != nil {
		return nil, err
	}

	x.redis = redis.NewClient(opt)

	return x, nil
}

func (r *RedisStore) Put(k string, v interface{}) error {

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := r.redis.Set(k, b, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) Get(k string, v interface{}) error {
	if b, err := r.redis.Get(k).Bytes(); err == nil {
		return json.Unmarshal(b, v)
	}

	return fmt.Errorf("no key")
}

func (r *RedisStore) Del(k string) error {
	if err := r.redis.Del(k).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) Keys(f string) ([]string, error) {
	keys, err := r.redis.Keys("*").Result()
	out := []string{}

	for _, k := range keys {
		if m, _ := path.Match(f, k); m {
			out = append(out, k)
		}
	}

	return out, err
}

func (r *RedisStore) Exists(k string) (bool, error) {
	i, err := r.redis.Exists(k).Result()
	return i == 1, err
}

func (r *RedisStore) Close() error { return r.redis.Close() }
