package conf

import (
	"fmt"
	api "github.com/hashicorp/consul/api"
	"log"
)

type consulSource struct {
	c         *api.Client
	namespace string // namespace is the prefix of the key, like "quickio", normally it is the name of the project
}

func NewConsulSource(namespace, addr string) Source {
	return consulSource{
		c:         NewConsulClient(addr),
		namespace: namespace,
	}
}

func (consul consulSource) ReadConf(env ENVType) *Config {
	pair, _, err := consul.c.KV().Get(consul.getConfKey(env), nil)
	if err != nil {
		panic("[Consul] Read Conf Error: " + err.Error() + ", key: " + consul.getConfKey(env))
	}
	if pair == nil {
		panic("[Consul] Read Conf Error: " + "key not found, key: " + consul.getConfKey(env))
	}
	fmt.Println("[Consul] Read Conf: \n", string(pair.Value))
	conf, err := UnmarshalConf(string(pair.Value))
	return conf
}

func (consul consulSource) getConfKey(envType ENVType) string {
	return consul.namespace + "/" + string(envType)
}

func (consul consulSource) ListenConf(env ENVType, onChange func(config *Config)) {
	panic("[Consul] unsupport listen conf")
}

func NewConsulClient(addr string) *api.Client {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	return client
}
