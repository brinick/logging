package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Logger defines the interface for logging clients
type Logger interface {
	Name() string
	Configurer
	LogLeveler
}

// Configurer defines the interface to configure logging clients
type Configurer interface {
	Configure(*Config)
}

// LogLeveler defines the interface for log level methods
type LogLeveler interface {
	Debug(string, ...Field)
	Info(string, ...Field)
	Error(string, ...Field)
	Fatal(string, ...Field)
}

// Config is the concrete type that is passed to a Configurer
type Config struct {
	LogLevel  string // Debug | Info | Error
	OutFormat string // json | text
	Outfile   string // path to file. Missing = send to stdout/err
}

// F is a shortcut for creating logging Fields
func F(name string, val interface{}) Field {
	return Field{
		Name: name,
		Val:  val,
	}
}

// Field represents a logging Field
type Field struct {
	Name string
	Val  interface{}
}

var (
	// logger is the logging package log client set via the SetClient function
	logger Logger = NewNullLogger()

	// ErrField is a shortcut function for adding an error field to the log output
	ErrField = func(e error) Field {
		return F("err", e)
	}
)

// NewClient returns a new instance of the concrete
// logging Client with the given name
func NewClient(name string) Logger {
	var logger Logger
	switch name {
	case "logrus":
		logger = NewLogrusLogger()
	default:
		logger = NewNullLogger()
	}

	return logger
}

// SetClient is a factory function to initiate the logging client
// with the given name. The instance is then set at the package level,
// and is retrievable in other packages using the Client() function.
func SetClient(name string) {
	if logger != nil && logger.Name() == name {
		return
	}
	switch name {
	case "logrus":
		logger = NewLogrusLogger()
	default:
		logger = NewNullLogger()
	}
}

// Client returns the logging client, or nil if it has not
// been initiated yet.
func Client() Logger {
	return logger
}

// Configure will configure the logger with the given attributes
func Configure(level, format, outfile string) {
	logger.Configure(&Config{
		LogLevel:  level,
		OutFormat: format,
		Outfile:   outfile,
	})
}

// ------------------------------------------------------------------
// Short cuts to the logging client
// ------------------------------------------------------------------

// Debug calls the logger Debug method
func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

// Info calls the logger Info method
func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

// Error calls the logger Error method
func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

// Fatal calls the logger Fatal method
func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

// ------------------------------------------------------------------

// source will return the line/lineno that called the
// given logging level function
func source() []Field {
	// source will return the line/lineno that caulled the
	// given logging level function
	// Who called the logging function.
	// As this function is called from the logger,
	// we need to go up 2 frames to get to the caller
	// of the logger function
	var (
		pkg = "???"
		src = "???:0"
	)

	pc, _, lineno, ok := runtime.Caller(3)
	if ok {
		caller := runtime.FuncForPC(pc).Name()
		path := filepath.Dir(caller)
		base := filepath.Base(caller)
		srcToks := strings.SplitN(base, ".", 2)

		pkg = filepath.Join(path, srcToks[0])
		src = fmt.Sprintf("%s:%d", srcToks[1], lineno)
	}
	return []Field{
		Field{"pkg", pkg},
		Field{"src", src},
	}
}
