package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// ------------------------------------------------------------------

func defaultLogrusConfig() *Config {
	return &Config{
		OutFormat: "text",
		LogLevel:  "info",
	}
}

// NewLogrusLogger wraps a logrus client
func NewLogrusLogger(cfg *Config) Logger {
	l := &LogrusLogger{
		log: logrus.New(),
	}
	if cfg == nil {
		cfg = defaultLogrusConfig()
	}

	l.Configure(defaultLogrusConfig().Update(cfg))
	return l
}

// ------------------------------------------------------------------

// LogrusLogger defines a logger using the logrus package as its backend
type LogrusLogger struct {
	log  *logrus.Logger
	path string
}

// Name returns the name of the logg
func (l *LogrusLogger) Name() string {
	return "logrus"
}

// Path returns the full path to the logger output, or empty string if
// logging is not going to a file
func (l *LogrusLogger) Path() string {
	return l.path
}

// Configure permits configuration of the logger via a Config struct
func (l *LogrusLogger) Configure(cfg *Config) {
	l.log.Out = os.Stdout

	l.path = cfg.Outfile
	if cfg.Outfile != "" {
		file, err := os.OpenFile(cfg.Outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			l.Error(
				"Failed to open file for logging, log to stdout/err instead",
				F("err", err),
				F("file", cfg.Outfile),
			)

			l.path = ""
			file = os.Stdout
		}

		l.log.Out = file
	}

	l.log.Level = l.toLogLevel(cfg.LogLevel)
	l.log.Formatter = l.toOutputFormat(cfg.OutFormat)
}

// Debug defines the debug level for this logger
func (l *LogrusLogger) Debug(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Debug(msg)
}

// Info defines the info level for this logger
func (l *LogrusLogger) Info(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Info(msg)
}

// Error defines the error level for this logger
func (l *LogrusLogger) Error(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Error(msg)
}

// Fatal defines the fatal level for this logger
func (l *LogrusLogger) Fatal(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Fatal(msg)
}

func (l *LogrusLogger) toOutputFormat(name string) logrus.Formatter {
	var formatter logrus.Formatter

	switch name {
	case "json":
		formatter = &logrus.JSONFormatter{}
	case "text":
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	default:
		panic("Unknown formatter " + name + ". Legal: json | text")
	}

	return formatter
}

func (l *LogrusLogger) toLogLevel(name string) logrus.Level {
	name = strings.TrimSpace(name)

	switch name {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	default:
		var msg = fmt.Sprintf(
			"Unknown log level: %s. Legal values: debug, info, error",
			name,
		)

		if len(name) == 0 {
			msg = "Please provide a log level. Legal values: debug, info, error"
		}
		panic(msg)
	}
}

// ------------------------------------------------------------------

func mapify(fields ...Field) map[string]interface{} {
	data := map[string]interface{}{}
	for _, f := range fields {
		data[f.Name] = f.Val
	}

	return data
}
