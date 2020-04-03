package logging

//NewNullLogger creates a new NullLogger
func NewNullLogger(cfg *Config) (*NullLogger, error) {
	l := &NullLogger{}
	l.Configure(cfg)
	return l, nil
}

// NullLogger does nothing
type NullLogger struct{}

// Name returns the name of the logger
func (NullLogger) Name() string {
	return "null"
}

// Path returns the path to the logger output file
func (NullLogger) Path() string {
	return ""
}

// Configure permits configuration of the logger via a Config struct
func (NullLogger) Configure(*Config) error { return nil }

// Debug defines the debug level for this logger
func (NullLogger) Debug(string, ...Field) {}

// DebugL defines the debug level for this logger
func (NullLogger) DebugL([]string, ...Field) {}

// Info defines the info level for this logger
func (NullLogger) Info(string, ...Field) {}

// InfoL defines the info level for this logger
func (NullLogger) InfoL([]string, ...Field) {}

// Error defines the error level for this logger
func (NullLogger) Error(string, ...Field) {}

// ErrorL defines the error level for this logger
func (NullLogger) ErrorL([]string, ...Field) {}

// Fatal defines the fatal level for this logger
func (NullLogger) Fatal(string, ...Field) {}

// FatalL defines the fatal level for this logger
func (NullLogger) FatalL([]string, ...Field) {}
