package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/TwiN/go-color"
)

// Logger represents a simple logger with different logging levels
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

// New creates a new Logger
func New() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime),
	}
}

// Info logs a message at info level
func (l *Logger) Info(msg string) {
	l.infoLogger.Printf(color.InBlue("%s\n"), msg)
}

// Error logs a message at error level
func (l *Logger) Error(msg string) {
	l.errorLogger.Printf(color.InRed("%s\n"), msg)
}

// Fatal logs a message at fatal level and exits the program
func (l *Logger) Fatal(msg string) {
	l.fatalLogger.Printf(color.InRed("%s\n"), msg)
}

func (l *Logger) PrintParserErrors(errors []string) {
	fmt.Println(color.InBold(color.InRed("parser errors:\n")))
	for _, msg := range errors {
		fmt.Println(color.InRed("\t" + msg))
	}
}
