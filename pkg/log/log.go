package log

type Logger interface {
	Warn(kv ...interface{})
	Err(kv ...interface{})
	Debug(kv ...interface{})
	Info(kv ...interface{})
	Printf(format string, v ...interface{})

	WithPrefix(kv ...string) Logger
}
