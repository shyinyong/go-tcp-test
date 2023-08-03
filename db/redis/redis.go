package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shyinyong/go-tcp-test/config"
	"time"
)

type RedisDB struct {
	client *redis.Client
}

func NewRedisDB(config *config.Config) (*RedisDB, error) {
	//address := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "", // No password
		DB:       0,  // Default DB
	})

	// Ping the server to check if the connection was successful
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisDB{client: client}, nil
}

func (r *RedisDB) Close() error {
	return r.client.Close()
}
