package redis

import (
	"github.com/redis/go-redis/v9"
)

//
//func (r *DB) Close() error {
//	return r.client.Close()
//}

var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
