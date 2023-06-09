package db

import (
	"fmt"
	"github.com/xixiwang12138/exermon/conf"
	"github.com/xixiwang12138/exermon/elog"
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

func NewRdsClient(cf conf.MySQLConfig, options ...Option) *RdsClient {
	c := &RdsClient{connect: &struct {
		UserName string
		Password string
		Address  string
		Port     int
		Database string
	}{UserName: cf.UserName, Password: cf.Password, Address: cf.Address, Port: cf.Port, Database: cf.DbName}}
	for _, option := range options {
		option(c)
	}
	c.Connect()
	return c
}

func (c *RdsClient) Connect() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.connect.UserName, c.connect.Password,
		c.connect.Address, c.connect.Port, c.connect.Database)
	if c.db, err = gorm.Open(mysql.Open(dsn)); err != nil { //TODO: 连接池等高级配置
		log.Fatal("open rds connection error: ", err.Error())
	}
	c.db.Logger = &elog.GormLogger
	return
}

func (c *RdsClient) Gorm() *gorm.DB {
	return c.db
}

func (c *RdsClient) AsyncRdsStruct(models ...any) {
	err := c.db.Migrator().AutoMigrate(models...)
	if err != nil {
		log.Fatal("AsyncRdsStruct failed: ", err)
	}
}
