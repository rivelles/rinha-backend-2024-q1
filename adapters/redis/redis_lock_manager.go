package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLockManager struct {
	client *redis.Client
}

func NewRedisLockManager() RedisLockManager {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return RedisLockManager{
		client: client,
	}
}

func (r RedisLockManager) Acquire(key string) error {
	currentLock, err := r.client.Get(context.Background(), key).Result()
	if err != nil && currentLock != "" {
		return fmt.Errorf("LOCK_ALREADY_ACQUIRED")
	}
	_, err = r.client.Set(context.Background(), key, true, 100*time.Second).Result()
	return err
}

func (r RedisLockManager) Release(key string) error {
	_, err := r.client.Del(context.Background(), key).Result()
	return err
}
