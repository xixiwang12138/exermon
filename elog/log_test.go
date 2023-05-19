package elog

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	SetConfig(DEBUG, "./log")
	cl := WithContext(context.WithValue(context.Background(), TraceIdHeader, "test01"))
	cl.Info("init successfully: %s", "hhh")
}
