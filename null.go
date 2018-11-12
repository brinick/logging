package logging

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

// NullLogger does nothing
type NullLogger struct{}

func (NullLogger) Name() string {
	return "null"
}
func (NullLogger) Configure(*Config) {}

func (NullLogger) Debug(string, ...Field) {}
func (NullLogger) Info(string, ...Field)  {}
func (NullLogger) Error(string, ...Field) {}
func (NullLogger) Fatal(string, ...Field) {}
