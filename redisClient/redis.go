package redisClient


import (
	"fmt"

	"github.com/go-redis/redis/v7"

	"github.com/buyubaya/go-blog/config"
)


func Initialize(config *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password, // no password set
		DB:       0,  // use default DB
	})
	client.FlushAll()


	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Redis Client error", err)
	}
	fmt.Println("Redis Client ready", pong)
	

	return client;
}