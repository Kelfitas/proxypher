package main

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

func (l *Logger) Log(pattern string, params... interface{}) {
	pattern = l.prefix + pattern
	logger.Printf(pattern, params...)
}
