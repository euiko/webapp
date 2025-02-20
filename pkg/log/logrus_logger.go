package log

import (
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logrus *logrus.Logger
}

func (l *LogrusLogger) Log(level Level, msg *Log) {
	entry := l.logrus.
		WithFields(logrus.Fields(msg.fields)).
		WithTime(msg.ts)

	// only add error when is not nil
	if msg.err != nil {
		entry = entry.WithError(msg.err)
	}

	entry.Log(toLogrusLevel(level), msg.message)
}

func NewLogrusLogger(level Level) *LogrusLogger {
	l := logrus.New()
	l.SetLevel(toLogrusLevel(level))
	return &LogrusLogger{logrus: l}
}

func toLogrusLevel(level Level) logrus.Level {
	switch level {
	case TraceLevel:
		return logrus.TraceLevel
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarningLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
