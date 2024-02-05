package id

import "testing"

func TestGenerate(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Log(Generate())
	}
}
