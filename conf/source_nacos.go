package conf

import (
	"log"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosSource struct {
	namespace    string
	group        string
	configClient config_client.IConfigClient
}

func (c *nacosSource) ListenConf(env ENVType, onChange func(config *Config)) {
	c.configClient.ListenConfig(vo.ConfigParam{
		DataId: string(env),
		Group:  c.namespace,
		OnChange: func(namespace, group, dataId, data string) {
			if string(env) != dataId || group != c.group || namespace != c.namespace {
				return
			}
			log.Println("nacos config changed: " + data)
			conf, err := UnmarshalConf(data)
			if err != nil {
				panic(err)
			}
			onChange(conf)
		},
	})
}

func (c *nacosSource) ReadConf(env ENVType) *Config {
	str, err := c.configClient.GetConfig(vo.ConfigParam{
		DataId: string(env),
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic("Read Conf Error: " + err.Error())
	}
	conf, err := UnmarshalConf(str)
	if err != nil {
		panic(err)
	}
	return conf
}

func NewNacosSource(namespace string, nacosServer string) Source {
	c := &nacosSource{
		namespace:    namespace,
		group:        "DEFAULT_GROUP",
		configClient: nil,
	}
	// 将nacosServer字符串解析为ip和port
	nacosServer = strings.TrimSpace(nacosServer)
	nacosServer = strings.Trim(nacosServer, "http://")
	nacosServer = strings.Trim(nacosServer, "https://")
	nacosServer = strings.Trim(nacosServer, "/")

	// 将nacosServer安装冒号分割为ip和port
	nacosServerArr := strings.Split(nacosServer, ":")
	if len(nacosServerArr) != 2 {
		panic("nacos server config error: " + nacosServer)
	}
	serverAddress := nacosServerArr[0]
	serverPortUint64, _ := strconv.ParseUint(nacosServerArr[1], 10, 64)

	sc := []constant.ServerConfig{{
		IpAddr: serverAddress,
		Port:   serverPortUint64,
	}}
	cc := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		LogLevel:            "error",
	}

	//基于nacos配置中心客户端
	c.configClient, _ = clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	return c
}
