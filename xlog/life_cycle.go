package xlog

import (
	"os"

	"github.com/spf13/viper"

	"github.com/xixiwang12138/exermon/component"
)

var LoggerComponent = &BaseLogger{}

func (b *BaseLogger) ConfigKeys() []string {
	return []string{"log.level", "log.path"}
}

func (b *BaseLogger) Name() string {
	return "log"
}

func (b *BaseLogger) Start(vp *viper.Viper) {
	vp.SetDefault("log.level", string(logLevelToName[INFO]))
	b.level = logNameToLevel[vp.GetString("log.level")]

	if !vp.IsSet("log.path") {
		b.file = os.Stdout
		return
	}
	out, err := os.Open(vp.GetString("log.path"))
	component.RaiseIfComponentError(err, "open log file")
	b.file = out
}

func (b *BaseLogger) Close() {
	b.file.Close()
}
