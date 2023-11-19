package elog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindPercentIndices(t *testing.T) {
	indices := findPercentIndices("%s abc %j %d")
	assert.Equal(t, len(indices), 3)
}

func TestReplaceJsonHolder(t *testing.T) {
	format := "%j abc %j %+d%j %+v"
	args := []interface{}{map[string]string{"abc": "123"}, []int{1, 2, 4}, 3, 0}

	newFormat, newArgs := replaceJsonHolder(format, args)
	assert.Equal(t, newFormat, `{"abc":"123"} abc [1,2,4] %+d0 %+v`)
	assert.Equal(t, len(newArgs), 1)
}
