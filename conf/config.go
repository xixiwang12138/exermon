package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"log"
	"strings"
)

// env 配置

func Global() *Config {
	return Context.global
}

func ENV() ENVType {
	return Context.env
}

type ENVType string

func (e ENVType) Test() bool {
	if e == TEST || e == DEV {
		return true
	}
	return false
}

const (
	PROD ENVType = "prod"
	DEV  ENVType = "dev"
	TEST ENVType = "test"
)

// MySQLConfig MySQL配置信息
type MySQLConfig struct {
	UserName string `toml:"user_name"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	Port     string `toml:"port"`
	DbName   string `toml:"db_name"`
}

// RedisConfig Redis配置信息，注意字段名与Options保持一致
type RedisConfig struct {
	Addr     string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// Config 所有配置信息]
type Config struct {
	HTTP       *BaseApiConfig    `toml:"http"`
	Mysql      *MySQLConfig      `toml:"mysql"`
	Redis      *RedisConfig      `toml:"redis"`
	SSL        *SSLConfig        `toml:"ssl"`
	FileConfig *FileUploadConfig `toml:"file"`
}

func (c *Config) check() {
	c.HTTP.check()
}

type FileUploadConfig struct {
	UploadPath      string `toml:"upload_path"`
	ImageAccessPath string `toml:"image_access_path"`
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

var Context = &context{}

type context struct {
	global *Config
	env    ENVType
}

func (c *context) Setup(buffer []byte, env ENVType) {
	c.env = env
	c.global = new(Config)
	err := toml.Unmarshal(buffer, c.global)
	if err != nil {
		log.Fatalf("[Get Config]Unmarshal toml: %v", err)
	}
}
