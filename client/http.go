package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/xixiwang12138/exermon/elog"
	"net"
	"net/http"
	"time"
)

var (
	DefaultHeader http.Header = nil
	NoBody        *any        = nil
)

type Option func(client *HttpRPCClient)

func WithHeader(header http.Header) Option {
	return func(client *HttpRPCClient) {
		client.header = header
	}
}

func WithRetry(getRetry, postRetry int) Option {
	return func(client *HttpRPCClient) {
		client.getRetry = getRetry
		client.postRetry = postRetry
	}
}

type HttpRPCClient struct {
	*http.Client
	postRetry int
	getRetry  int

	header http.Header
}

func NewHttpRPCClient(options ...Option) *HttpRPCClient {
	c := &HttpRPCClient{
		postRetry: 0,
		getRetry:  5,
	}
	for _, option := range options {
		option(c)
	}
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

func (h *HttpRPCClient) JsonPost(ctx context.Context, url string, header http.Header, data any, result any) error {
	request, err := h.getJSONRequest(ctx, "POST", url, header, data)
	if err != nil {
		return err
	}
	return h.do(ctx, request, h.postRetry, result)
}

func (h *HttpRPCClient) JsonGet(ctx context.Context, url string, header http.Header, result any) error {
	request, err := h.getJSONRequest(ctx, "GET", url, header, NoBody)
	if err != nil {
		return err
	}
	return h.do(ctx, request, h.getRetry, result)
}

func (h *HttpRPCClient) getJSONRequest(ctx context.Context, method string, url string, header http.Header, data any) (request *http.Request, err error) {
	el := elog.WithContext(ctx)
	el.Info("begin request: %s %s", method, url)
	var (
		buf []byte
	)
	if data != NoBody {
		if buf, err = json.Marshal(data); err != nil {
			el.Error("marshal request data error: %s", err.Error())
			return nil, err
		}
	}

	request, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(buf))
	if err != nil {
		el.Error("new request error: %s", err.Error())
		return nil, err
	}

	if header == nil {
		request.Header = h.header
	} else {
		request.Header = header
	}

	return
}

func (h *HttpRPCClient) do(ctx context.Context, request *http.Request, tryTimes int, result any) error {
	el := elog.WithContext(ctx)
	var (
		err  error
		resp *http.Response
	)
	for tryTimes >= 0 {
		resp, err = h.Do(request)
		if err != nil {
			el.Error("post error: %s, tryTimes left: %d", err.Error(), tryTimes)
			goto sleep
		}
		if resp.StatusCode/100 == 2 {
			if result == NoBody {
				return nil
			}
			err = json.NewDecoder(resp.Body).Decode(result)
			if err != nil {
				el.Error("unmarshal response error: %s", err.Error())
				return err
			}
			return nil
		}
		if resp.StatusCode/100 == 4 {
			el.Error("client error, no need retry, status: %s", resp.Status)
			return errors.New(resp.Status)
		}
	sleep:
		tryTimes--
		el.Warning("do request error, status: %s, begin retry", resp.Status)
		time.Sleep(time.Duration(h.postRetry-tryTimes) * time.Second)
	}
	return err
}
