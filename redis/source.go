package redis

import (
	"context"
	"crypto/tls"
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
	Address    string
	Password   string
	DB         int
	TLS        bool
	SNI        string
	SkipVerify bool
}

type RedisClient struct {
	callee callee
	Client *redis.Client

	extendedCallees []*callee
	ExtendedClients []*redis.Client
}

func Setup(cf *conf.RedisConfig, extendedCf []*conf.RedisConfig, options ...Option) {
	Component = &RedisClient{
		callee: callee{Address: cf.Addr, Password: cf.Password, DB: cf.DB, TLS: cf.TLS, SNI: cf.SNI}}

	for _, extCallee := range extendedCf {
		Component.extendedCallees = append(Component.extendedCallees, &callee{
			Address:  extCallee.Addr,
			Password: extCallee.Password,
			DB:       extCallee.DB,
			TLS:      extCallee.TLS,
			SNI:      extCallee.SNI,
		})
	}

	for _, option := range options {
		option(Component)
	}
	Component.Connect()
}

func (c *RedisClient) newRedisClient(callee *callee) *redis.Client {
	var tlsConfig *tls.Config = nil
	if callee.TLS {
		tlsConfig = &tls.Config{
			ServerName:         callee.SNI,
			InsecureSkipVerify: callee.SkipVerify,
		}
	}

	return redis.NewClient(&redis.Options{
		Addr:      callee.Address,
		Password:  callee.Password,
		DB:        callee.DB,
		TLSConfig: tlsConfig,
	})
}

func (c *RedisClient) testClient(client *redis.Client) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	r, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("connect redis failed: ", r, err.Error())
	}
	log.Printf("connect redis %s db%d successfully!!", client.Options().Addr, client.Options().DB)
}

func (c *RedisClient) Connect() {
	c.Client = c.newRedisClient(&c.callee)
	c.testClient(c.Client)

	for _, callee := range c.extendedCallees {
		client := c.newRedisClient(callee)
		c.ExtendedClients = append(c.ExtendedClients, client)
		c.testClient(client)
	}
	return
}
