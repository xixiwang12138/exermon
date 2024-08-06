package conf

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type ENVType string

const (
	PROD ENVType = "prod"
	DEV  ENVType = "dev"
	TEST ENVType = "test"
)

// MySQLConfig MySQL配置信息
type MySQLConfig struct {
	UserName string `toml:"user_name" yaml:"user_name" json:"user_name"`
	Password string `toml:"password" yaml:"password" json:"password"`
	Address  string `toml:"address" yaml:"address" json:"address"`
	Port     int    `toml:"port" yaml:"port" json:"port"`
	DbName   string `toml:"db_name" yaml:"db_name" json:"db_name"`
}

// RedisConfig Redis配置信息
type RedisConfig struct {
	Addr     string `toml:"address" yaml:"addr" json:"addr"`
	Password string `toml:"password" yaml:"password" json:"password"`
	DB       int    `toml:"db" yaml:"db" json:"db"`
}

// LogConfig 日志配置信息
type LogConfig struct {
	Level   string `toml:"level" yaml:"level" json:"level"`
	DirPath string `toml:"dir_path" yaml:"dir_path" json:"dir_path"`
}

// Config 所有配置信息
type Config struct {
	HTTP       *BaseApiConfig    `toml:"http" yaml:"http" json:"http"`
	Mysql      *MySQLConfig      `toml:"mysql" yaml:"mysql" json:"mysql"`
	Redis      *RedisConfig      `toml:"redis" yaml:"redis" json:"redis"`
	ExtendedRedis []*RedisConfig `toml:"extended_redis" yaml:"extended_redis" json:"extended_redis"`
	SSL        *SSLConfig        `toml:"ssl" yaml:"ssl" json:"ssl"`
	FileConfig *FileUploadConfig `toml:"file" yaml:"file" json:"file"`
	LogConfig  *LogConfig        `toml:"log" yaml:"log" json:"log"`
}

func (c *Config) check() {
	c.HTTP.check()
}

type FileUploadConfig struct {
	UploadPath      string `toml:"upload_path" yaml:"upload_path" json:"upload_path"`
	ImageAccessPath string `toml:"image_access_path" yaml:"image_access_path" json:"image_access_path"`
}

type SSLConfig struct {
	KeyPath  string `toml:"key"`
	CertPath string `toml:"cert"`
}

// BaseApiConfig web层配置信息
type BaseApiConfig struct {
	PORT  string `toml:"port"`
	group map[string]func(*gin.RouterGroup)
}

func (c BaseApiConfig) check() {
	if !strings.HasPrefix(c.PORT, ":") {
		log.Fatalf("http listen port should start with : ,but now is %s \n", c.PORT)
	}
}
