package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Option func(*RdsClient) *RdsClient

var (
	WithConnPoolSize = func(size int) Option {
		return func(client *RdsClient) *RdsClient {
			client.pool = size
			return client
		}
	}
)

type RdsClient struct {
	db *gorm.DB

	connect *struct {
		UserName string
		Password string
		Address  string
		Port     int
		Database string
	}

	pool int
}

func NewRdsClient(port int, userName, password, address, database string, options ...Option) *RdsClient {
	c := &RdsClient{connect: &struct {
		UserName string
		Password string
		Address  string
		Port     int
		Database string
	}{UserName: userName, Password: password, Address: address, Port: port, Database: database}}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *RdsClient) Connect() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.connect.UserName, c.connect.Password,
		c.connect.Address, c.connect.Port, c.connect.Database)
	if c.db, err = gorm.Open(mysql.Open(dsn)); err != nil { //TODO: 连接池等高级配置
		log.Fatal("open rds connection error: ", err.Error())
	}
	return nil
}

func (c *RdsClient) Gorm() *gorm.DB {
	return c.db
}
