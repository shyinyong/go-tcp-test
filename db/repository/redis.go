package repository

import "github.com/shyinyong/go-tcp-test/db/redis"

type RedisRepository struct {
	db *redis.RedisDB
}

func NewRedisRepository(db *redis.RedisDB) *RedisRepository {
	return &RedisRepository{db: db}
}
