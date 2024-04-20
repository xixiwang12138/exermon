package conf

// Source is the interface that wraps the basic confSource.
type Source[T any] interface {
	ReadConf(env ENVType) *T
	ListenConf(env ENVType, onChange func(config *T))
}
