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

type Configurer interface {
	Configure(*Config)
}

type LogLeveler interface {
	Debug(string, ...Field)
	Info(string, ...Field)
	Error(string, ...Field)
	Fatal(string, ...Field)
}

type Config struct {
	LogLevel  string
	OutFormat string
	Outfile   string
}

// Shortcut function for creating logging Fields
func F(name string, val interface{}) Field {
	return Field{
		Name: name,
		Val:  val,
	}
}

type Field struct {
	Name string
	Val  interface{}
}

// logger is the logging package log client set via the SetClient function
var logger Logger = NewNullLogger()

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

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

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
