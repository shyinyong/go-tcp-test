package store

import "context"
import "github.com/redis/go-redis/v9"

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(address string, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

func (c *RedisClient) Set(ctx context.Context, key string, value string) error {
	err := c.client.Set(ctx, key, value)
	// Implement the Set operation
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	// Implement the Get operation
}
