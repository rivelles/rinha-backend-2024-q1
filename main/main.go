package main

import (
	"github.com/rivelles/rinha-backend-2024-q1/adapters/redis"
	"github.com/rivelles/rinha-backend-2024-q1/adapters/scylla"
)

func main() {
	repository := scylla.NewScyllaRepository()
	lockManager := redis.NewRedisLockManager()
	app := NewApp(repository, lockManager)

	app.Run(9091)
}
