package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type HttpRPCClient struct {
	*http.Client
}

func NewHttpRPCClient() *HttpRPCClient {
	c := &HttpRPCClient{}
	c.build()
	return c
}

func (h *HttpRPCClient) build() {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// InsecureSkipVerify用来控制客户端是否证书和服务器主机名。
		// 如果设置为true, 则不会校验证书以及证书中的主机名和服务器主机名是否一致
		// 由于用于RPC，一般认为配置的主机内网安全
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, //拨号等待连接完成的最大时间
			KeepAlive: 30 * time.Minute, //保持网络连接活跃keep-alive探测间隔时间
		}).DialContext,
		MaxIdleConns:        200,
		IdleConnTimeout:     300 * time.Second,
		MaxIdleConnsPerHost: 20,
	}
	h.Client = &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second, //设置超时，包含connection时间、任意重定向时间、读取response body时间
	}
}
