package id

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Log(Generate())
	}
}

func TestGenerate2(t *testing.T) {
	os.Setenv("NODE_ID", "2")
	t.Log(Generate())
}
