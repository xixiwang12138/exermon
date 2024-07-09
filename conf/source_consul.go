package conf

import (
	"fmt"
	"log"

	api "github.com/hashicorp/consul/api"
)

type consulSource[T any] struct {
	c         *api.Client
	namespace string // namespace is the prefix of the key, like "quickio", normally it is the name of the project
}

func NewConsulSource[T any](namespace, addr string, token string) Source[T] {
	return consulSource[T]{
		c:         NewConsulClient(addr, token),
		namespace: namespace,
	}
}

func (consul consulSource[T]) ReadConf(env ENVType) *T {
	pair, _, err := consul.c.KV().Get(consul.getConfKey(env), nil)
	if err != nil {
		panic("[Consul] Read Conf Error: " + err.Error() + ", key: " + consul.getConfKey(env))
	}
	if pair == nil {
		panic("[Consul] Read Conf Error: " + "key not found, key: " + consul.getConfKey(env))
	}
	fmt.Println("[Consul] Read Conf: \n", string(pair.Value))
	viperParser(string(pair.Value))
	conf, err := UnmarshalConf[T](string(pair.Value))
	return conf
}

func (consul consulSource[T]) getConfKey(envType ENVType) string {
	return consul.namespace + "/" + string(envType)
}

func (consul consulSource[T]) ListenConf(env ENVType, onChange func(config *T)) {
	panic("[Consul] unsupport listen conf")
}

func NewConsulClient(addr string, token string) *api.Client {
	config := api.DefaultConfig()
	config.Address = addr
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	return client
}
