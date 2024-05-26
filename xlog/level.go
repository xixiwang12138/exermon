package xlog

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	SILENT = 100
)

var logLevelToName = map[LogLevel][]byte{
	DEBUG: []byte("DEBUG"),
	INFO:  []byte("INFO"),
	WARN:  []byte("WARN"),
	ERROR: []byte("ERROR"),
	FATAL: []byte("FATAL"),
}

var logNameToLevel = map[string]LogLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

const (
	left  byte = '['
	right byte = ']'
	black byte = ' '
	colon byte = ':'
)

var (
	rightAndLeft []byte = []byte("][")
)
