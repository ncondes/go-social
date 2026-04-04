package logging

import "go.uber.org/zap"

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(logger *zap.SugaredLogger) Logger {
	return &zapLogger{logger: logger}
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...any) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
