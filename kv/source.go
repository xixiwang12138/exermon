package kv

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/xixiwang12138/exermon/conf"
	"log"
	"time"
)

type Option func(client *RedisClient) *RedisClient

type RedisClient struct {
	connect struct {
		Address  string
		Password string
		DB       int
	}
	Client *redis.Client
}

func NewRedisClient(cf conf.RedisConfig, options ...Option) *RedisClient {
	c := &RedisClient{
		connect: struct {
			Address  string
			Password string
			DB       int
		}{Address: cf.Addr, Password: cf.Password, DB: cf.DB}}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *RedisClient) Connect() {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     c.connect.Address,
		Password: c.connect.Password,
		DB:       c.connect.DB,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("connect redis failed: ", r, err.Error())
	}
	log.Println("connect redis successfully!")
	return
}
