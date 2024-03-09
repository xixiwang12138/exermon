package conf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	example := `
[ex]
data = 1
`
	Viper.SetConfigType("toml")
	Viper.ReadConfig(strings.NewReader(example))

	t.Log(Viper.AllKeys())
	res := Viper.GetInt("ex.data")
	assert.Equal(t, 1, res)
}
