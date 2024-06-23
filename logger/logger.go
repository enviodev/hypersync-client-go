package logger

// global is the package-level global logger instance.
var global Logger

// L returns the global logger instance.
func L() Logger {
	return global
}

// SetGlobalLogger sets the global logger instance to the provided Logger.
func SetGlobalLogger(l Logger) {
	global = l
}

// Logger is the interface for structured logging. It provides methods for logging
// at various levels of severity.
type Logger interface {
	// Debug logs a message at the debug level.
	Debug(format string, args ...interface{})
	// Info logs a message at the info level.
	Info(format string, args ...interface{})
	// Warn logs a message at the warn level.
	Warn(format string, args ...interface{})
	// Error logs a message at the error level.
	Error(format string, args ...interface{})
	// Fatal logs a message at the fatal level. The application will terminate after logging the message.
	Fatal(format string, args ...interface{})
	// Panic logs a message at the panic level. The application will panic after logging the message.
	Panic(format string, args ...interface{})
}
