package logging

type Logger interface {
	Errorw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
	Sync() error
}
