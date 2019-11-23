package dal

import (
    "fmt"
    "github.com/go-redis/redis"
    "sync"
    "time"
)

var KVs KVStore

type KVStore interface {
    Set(k, v string, ttl time.Duration) error
    Get(k string) (string, error)
}

type MemBackend struct {
    sync.RWMutex
    store map[string]string
}

func NewMemBackend() *MemBackend {
    return &MemBackend{store: make(map[string]string)}
}

func (s *MemBackend) Get(key string) (string, error) {
    s.RLock()
    defer s.RUnlock()
    value, ok := s.store[key]
    if !ok {
        return value, fmt.Errorf("get error")
    }
    return value, nil
}

func (s *MemBackend) Set(k, v string, ttl time.Duration) error {
    s.Lock()
    defer s.Unlock()
    s.store[k] = v
    return nil
}

type RedisBackend struct {
    redisC *redis.Client
}

func (r *RedisBackend) Set(k, v string, ttl time.Duration) error {
    cmd := r.redisC.Set(k, v, ttl)
    if err := cmd.Err(); err != nil {
        return err
    }
    return nil
}

func (r *RedisBackend) Get(k string) (string, error) {
    cmd := r.redisC.Get(k)
    value, err := cmd.Result()
    if err != nil {
        return "", err
    }
    return value, nil
}

func InitKV(redisAddr, MemoryOrRedis string) error {
    var store KVStore
    if MemoryOrRedis == "memory" {
        store = NewMemBackend()
    } else {
        redisC := redis.NewClient(&redis.Options{
            Addr:     redisAddr,
            Password: "",
            DB:       0, // use default DB
        })

        store = &RedisBackend{
            redisC: redisC,
        }
    }
    KVs = store
    return nil
}
