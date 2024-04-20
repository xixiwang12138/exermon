package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/xixiwang12138/exermon/conf"
	"log"
	"time"
)

var Component *RedisClient

const (
	Nil = redis.Nil
)

type Option func(client *RedisClient) *RedisClient

type RedisClient struct {
	callee struct {
		Address  string
		Password string
		DB       int
	}
	Client *redis.Client
}

func Setup(cf *conf.RedisConfig, options ...Option) {
	Component = &RedisClient{
		callee: struct {
			Address  string
			Password string
			DB       int
		}{Address: cf.Addr, Password: cf.Password, DB: cf.DB}}
	for _, option := range options {
		option(Component)
	}
	Component.Connect()
}

func (c *RedisClient) Connect() {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     c.callee.Address,
		Password: c.callee.Password,
		DB:       c.callee.DB,
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)
	r, err := c.Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("connect redis failed: ", r, err.Error())
	}
	log.Println("connect redis successfully!!")
	return
}
