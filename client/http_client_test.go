package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xixiwang12138/exermon/elog"
	"testing"
)

func TestName(t *testing.T) {
	c := NewHttpRPCClient()
	d := make(map[string]any)
	elog.SetConfig(elog.DEBUG, "./log")
	c.JsonGet(context.WithValue(context.Background(), elog.TraceIdHeader, "test"), "http://43.139.80.71:8081/api//", DefaultHeader, &d)

	v, _ := json.Marshal(d)
	fmt.Println(string(v))

}

func TestName2(t *testing.T) {
	fmt.Println(len("[test][23:29:36][DEBUG] http.go:144"))
}