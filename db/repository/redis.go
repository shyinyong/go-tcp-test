package repository

import "github.com/shyinyong/go-tcp-test/db/redis"

type RedisRepository struct {
	db *redis.DB
}

func NewRedisRepository(db *redis.DB) *RedisRepository {
	return &RedisRepository{db: db}
}
