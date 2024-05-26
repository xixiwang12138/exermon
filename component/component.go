package component

import "github.com/spf13/viper"

type Component interface {
	Name() string
	ConfigKeys() []string
	Start(conf *viper.Viper)
	Close()
}
