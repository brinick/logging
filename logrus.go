package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

// ------------------------------------------------------------------

// LogrusLogger wraps a logrus client
func NewLogrusLogger() Logger {
	return &LogrusLogger{
		log: logrus.New(),
	}
}

// ------------------------------------------------------------------

type LogrusLogger struct {
	log *logrus.Logger
}

func (l *LogrusLogger) Name() string {
	return "logrus"
}

func (l *LogrusLogger) Configure(cfg *Config) {
	l.log.Out = os.Stdout
	l.log.Level = l.toLogLevel(cfg.LogLevel)
	l.log.Formatter = l.toOutputFormat(cfg.OutFormat)
}

func (l *LogrusLogger) Debug(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Debug(msg)
}

func (l *LogrusLogger) Info(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Info(msg)
}

func (l *LogrusLogger) Error(msg string, fields ...Field) {
	fields = append(fields, source()...)
	l.log.WithFields(mapify(fields...)).Error(msg)
}

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
	switch name {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	default:
		panic("Unknown log level " + name + ". Legal: debug | info | error")
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
