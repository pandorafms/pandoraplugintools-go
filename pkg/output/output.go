package output

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
	Debug  bool
}

type LoggerOption func(*Logger)

func WithStdout(w io.Writer) LoggerOption {
	return func(l *Logger) { l.Stdout = w }
}

func WithStderr(w io.Writer) LoggerOption {
	return func(l *Logger) { l.Stderr = w }
}

func WithDebug(enabled bool) LoggerOption {
	return func(l *Logger) { l.Debug = enabled }
}

func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Debug:  false,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Logger) PrintStdout(format string, args ...interface{}) {
	fmt.Fprintf(l.Stdout, format+"\n", args...)
}

func (l *Logger) PrintStderr(format string, args ...interface{}) {
	fmt.Fprintf(l.Stderr, format+"\n", args...)
}

func (l *Logger) PrintDebug(format string, args ...interface{}) {
	if l.Debug {
		fmt.Fprintf(l.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

var defaultLogger = NewLogger()

func PrintStdout(format string, args ...interface{}) {
	defaultLogger.PrintStdout(format, args...)
}

func PrintStderr(format string, args ...interface{}) {
	defaultLogger.PrintStderr(format, args...)
}

func PrintDebug(format string, args ...interface{}) {
	defaultLogger.PrintDebug(format, args...)
}

func SetDebug(enabled bool) {
	defaultLogger.Debug = enabled
}
