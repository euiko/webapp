package log

import (
	"context"
	"time"
)

type Level int

const (
	FatalLevel Level = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type (
	// Logger represent any logging capable
	Logger interface {
		Log(level Level, msg *Log)
	}

	// Fields is a map of key-value pairs of metadata that will be logged
	Fields map[string]interface{}

	// Log holds the data that will be logged
	Log struct {
		ctx     context.Context
		ts      time.Time
		message string
		fields  Fields
		err     error
	}

	// Option is a function that configures a log
	Option interface {
		Configure(o *Log)
	}

	// OptionFunc is a function that configures a log
	OptionFunc func(o *Log)

	// contextKey is a type to help avoid context key collision
	contextKey int
)

// define context keys using iota
const (
	fieldsContextKey contextKey = iota
)

// globalLogger holds the default logger
var (
	globalLogger Logger
)

// Configure implements Option interface
func (f OptionFunc) Configure(o *Log) {
	f(o)
}

// WithFileds adds all parameter fields for decorate current message option
func WithFields(fields Fields) Option {
	return OptionFunc(func(log *Log) {
		for k, v := range fields {
			log.fields[k] = v
		}
	})
}

// WithField add a single field to the message option
func WithField(key string, value interface{}) Option {
	return OptionFunc(func(log *Log) {
		log.fields[key] = value
	})
}

// WithContext decorate current message option to include a context
// maybe override the logger option
func WithContext(ctx context.Context) Option {
	return OptionFunc(func(log *Log) {
		log.ctx = ctx
	})
}

// WithTime decorate message option that override default now timestamp
func WithTime(t time.Time) Option {
	return OptionFunc(func(log *Log) {
		log.ts = t
	})
}

// WithError decorate message option that override the error
func WithError(err error) Option {
	return OptionFunc(func(log *Log) {
		log.err = err
	})
}

// Fatal logs a message at the fatal level
func Fatal(msg string, opts ...Option) error {
	return log(FatalLevel, msg, opts...)
}

// Error logs a message at the error level
func Error(msg string, opts ...Option) error {
	return log(ErrorLevel, msg, opts...)
}

// Warning logs a message at the warning level
func Warning(msg string, opts ...Option) error {
	return log(WarningLevel, msg, opts...)
}

// Info logs a message at the info level
func Info(msg string, opts ...Option) error {
	return log(InfoLevel, msg, opts...)
}

// Debug logs a message at the debug level
func Debug(msg string, opts ...Option) error {
	return log(DebugLevel, msg, opts...)
}

// Trace logs a message at the trace level
func Trace(msg string, opts ...Option) error {
	return log(TraceLevel, msg, opts...)
}

// SetDefault sets the default logger
func SetDefault(logger Logger) {
	globalLogger = logger
}

// Default returns the default logger
func Default() Logger {
	return globalLogger
}

// SetFieldsContext sets the fields context
func SetFieldsContext(ctx context.Context, fields Fields) context.Context {
	// load or create fields from context
	current := getFieldsContext(ctx)
	if current == nil {
		current = make(Fields)
	}

	// merge all fields
	for k, v := range fields {
		current[k] = v
	}

	// inject to context
	return context.WithValue(ctx, fieldsContextKey, current)
}

func log(level Level, msg string, opts ...Option) error {
	// global logger not yet specified, then do nothing
	if globalLogger == nil {
		return nil
	}

	log := newLog(msg, opts...)
	// TODO: support logger other than global
	logger := globalLogger

	// load all context specific option
	if log.ctx != nil {
		ctx := log.ctx

		// add additional fields defined in context
		// override all previously defined option
		fields := getFieldsContext(ctx)
		for k, v := range fields {
			log.fields[k] = v
		}
	}

	// do log
	logger.Log(level, log)
	return nil
}

func newLog(message string, options ...Option) *Log {

	// instantiate log
	log := &Log{
		ts:      time.Now(),
		message: message,
		fields:  make(Fields),
		ctx:     nil,
		err:     nil,
	}

	// load all decorator function
	for _, o := range options {
		o.Configure(log)
	}

	return log
}

func getFieldsContext(ctx context.Context) Fields {
	instance := ctx.Value(fieldsContextKey)
	if instance == nil {
		return nil
	}

	fields, ok := instance.(Fields)
	if !ok {
		return nil
	}

	return fields
}

func init() {
	// set the default logger to use logrus with Info level
	SetDefault(NewLogrusLogger(InfoLevel))
}
