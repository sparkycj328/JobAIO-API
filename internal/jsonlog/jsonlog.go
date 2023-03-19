package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level defines a Level type to represent the severity level for a log entry.
type Level int8

// Initialize constants which represent a specific severity level.
// Usage of iota allows for automatic successive integers to be
// assigned to our constants
const (
	LevelInfo  Level = iota // Has the value of 0
	LevelError              // Has the value of 1
	LevelFatal              // Has the value of 2
	LevelOff                // Has the value of 3
)

// String returns a human-friendly string for the severity level
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelOff:
		return "OFF"
	default:
		return ""
	}
}

// Logger is a custom logger struct type. This holds the output destination that
// the log entries will be written for and a mutex which will coordinate the writes
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// PrintInfo will write errors at the Info level
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

// PrintError will write errors at the Error level
func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	// ensure that the level is not below the minimum severity
	// and then return with no further action.
	if level < l.minLevel {
		return 0, nil
	}
	// declare an anonymous struct that will hold the data for the log entry
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// Declare a line variable for holding the actual log entry text/
	var line []byte

	// Marshal the anonymous struct to JSON and store it in the line variable
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}
	// Lock the mutex so that no two writes to the output destination can happen concurrently.
	// If we don't do this, it's possible that the text for two or more log entries
	// will be intermingled in the output.
	l.mu.Lock()
	defer l.mu.Unlock()

	// Write the log entry followed by a newline
	return l.out.Write(append(line, '\n'))
}

// We also implement a Write() method on our Logger type so that it satisfies the
// io.Writer interface. This writes a log entry at the ERROR level with no additional
// properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
