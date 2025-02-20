package log

import (
	"strconv"
	"strings"
)

// ParseLevel parses a string level to a Level
func ParseLevel(level string) Level {
	// try parse as int
	parsed, err := strconv.ParseInt(level, 10, 64)
	if err == nil &&
		parsed >= int64(FatalLevel) &&
		parsed <= int64(TraceLevel) {
		return Level(parsed)
	}

	// try parse as string
	lowercased := strings.ToLower(level)
	switch lowercased {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarningLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		// use info as default
		return InfoLevel
	}
}
