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

type callee struct {
	Address  string
	Password string
	DB       int
}

type RedisClient struct {
	callee callee
	Client *redis.Client

	extenedCallees   []*callee
	ExtendedClients []*redis.Client
}

func Setup(cf *conf.RedisConfig, extendedCf []*conf.RedisConfig, options ...Option) {
	Component = &RedisClient{
		callee: struct {
			Address  string
			Password string
			DB       int
		}{Address: cf.Addr, Password: cf.Password, DB: cf.DB}}

	for _, extCallee := range extendedCf {
		Component.extenedCallees = append(Component.extenedCallees, &callee{
			Address:  extCallee.Addr,
			Password: extCallee.Password,
			DB:       extCallee.DB,
		})
	}

	for _, option := range options {
		option(Component)
	}
	Component.Connect()
}

func (c *RedisClient) newRedisClient(callee *callee) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     callee.Address,
		Password: callee.Password,
		DB:       callee.DB,
	})
}

func (c *RedisClient) testClient(client *redis.Client)  {
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)
	r, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("connect redis failed: ", r, err.Error())
	}
	log.Printf("connect redis %s db%d successfully!!", client.Options().Addr, client.Options().DB)
}

func (c *RedisClient) Connect() {
	c.Client = c.newRedisClient(&c.callee)
	c.testClient(c.Client)
	
	for _, callee := range c.extenedCallees {
		client := c.newRedisClient(callee)
		c.ExtendedClients = append(c.ExtendedClients, client)
		c.testClient(client)
	}
	return
}
