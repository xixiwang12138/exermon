package conf

// Source is the interface that wraps the basic confSource.
type Source interface {
	ReadConf(env ENVType) *Config
	ListenConf(env ENVType, onChange func(config *Config))
}
