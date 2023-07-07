package logger

import (
	"log"
	"os"
)

// Logger represents a simple logger with different logging levels
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

// NewLogger creates a new Logger
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime),
	}
}

// Info logs a message at info level
func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Error logs a message at error level
func (l *Logger) Error(msg string) {
	red := "\033[31m"
	reset := "\033[0m"
	l.errorLogger.Printf("%s%s%s\n", red, msg, reset)
}

// Fatal logs a message at fatal level and exits the program
func (l *Logger) Fatal(msg string) {
	red := "\033[31m"
	reset := "\033[0m"
	l.fatalLogger.Fatalf("%s%s%s\n", red, msg, reset)
}
