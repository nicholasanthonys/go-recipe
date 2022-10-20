package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

type redistHandler struct {
	client *redis.Client
}

func NewRedistHanlder(c *config) (*redistHandler, error) {
	_, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	db, err := strconv.Atoi(c.database)
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: c.password, // no password set
		DB:       db,         // use default DB
	})
	// pong, err := rdb.Ping(ctx).Result()
	// fmt.Println(pong, err)
	// if err != nil {
	// 	fmt.Println("error ping : ", err.Error())
	// 	log.Fatal(err)
	// }

	fmt.Println("Successfully connected to redis")

	return &redistHandler{
		client: rdb,
	}, nil
}

func (rh redistHandler) Set(ctx context.Context, key string, value string, time time.Duration) error {
	if err := rh.client.Set(ctx, key, value, time).Err(); err != nil {
		return err
	}
	return nil
}

func (rh redistHandler) Get(ctx context.Context, key string) (string, error) {
	return rh.client.Get(ctx, key).Result()
}

func (rh redistHandler) Del(ctx context.Context, key []string) (int64, error) {
	return rh.client.Del(ctx, key...).Result()
}
