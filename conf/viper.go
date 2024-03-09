package conf

import (
	"strings"

	"github.com/spf13/viper"
)

var Viper = viper.New()

func viperParser(confContent string) {
	err := Viper.ReadConfig(strings.NewReader(confContent))
	if err != nil {
		panic("read config: " + err.Error())
	}
}

func SetViperDefault(values map[string]any) {
	for k, v := range values {
		Viper.SetDefault(k, v)
	}
}
