package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/brinick/fs"
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
func (l *LogrusLogger) Configure(cfg *Config) error {
	l.log.Out = os.Stdout
	l.log.Level = l.toLogLevel(cfg.LogLevel)
	l.log.Formatter = l.toOutputFormat(cfg.OutFormat)

	l.path = strings.TrimSpace(cfg.Outfile)
	if l.path != "" {
		if err := l.logfileCheck(); err != nil {
			return err
		}

		file, err := os.OpenFile(cfg.Outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			return err
		}

		/*
				l.Error(
					"Failed to open file for logging, using stdout/err instead",
					F("err", err),
					F("file", cfg.Outfile),
				)

			}
		*/

		l.log.Out = file
	}

	return nil
}

// logfileCheck verifies, if logging to a file is requested, that the
// file parent directory exists
func (l *LogrusLogger) logfileCheck() error {
	logfile := fs.File(l.path)
	logfileDir := logfile.Dir()
	exists, err := logfileDir.Exists()
	if err != nil {
		return fmt.Errorf("unable to check if logfile")
	}

	if !exists {
		return fmt.Errorf(
			"log file parent directory inexistant, please create (%s)",
			logfileDir.Path,
		)
	}

	return nil
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
