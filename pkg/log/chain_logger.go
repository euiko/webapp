package log

type LoggerFactory func() Logger

type ChainLogger struct {
	loggers []Logger
}

func (l *ChainLogger) Log(level Level, msg *Log) {
	// simply range over through all available logger
	// TODO: support parallel log
	for _, l := range l.loggers {
		l.Log(level, msg)
	}
}

// NewChainLogger creates a new chained logger
// by supplying all LoggerFactory
func NewChainLogger(factories ...LoggerFactory) *ChainLogger {
	loggers := make([]Logger, len(factories))
	for i, f := range factories {
		loggers[i] = f()
	}

	return &ChainLogger{loggers: loggers}
}
