package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel represents the severity of the log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Color codes for console output
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
)

// Logger represents a logging instance
type Logger struct {
	MinLevel LogLevel
	ShowTime bool
}

// New creates a new Logger with default settings
func New() *Logger {
	return &Logger{
		MinLevel: INFO,
		ShowTime: true,
	}
}

// Log is the main logging function that handles different log levels
func (l *Logger) Log(level LogLevel, message string, args ...interface{}) {
	if level < l.MinLevel {
		return
	}

	// Format the message with additional args if provided
	formattedMsg := message
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(message, args...)
	}

	// Get current timestamp
	timestamp := ""
	if l.ShowTime {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	// Prepare level text and color
	var levelText, color string
	switch level {
	case DEBUG:
		levelText = "DEBUG"
		color = blue
	case INFO:
		levelText = "INFO"
		color = green
	case WARNING:
		levelText = "WARNING"
		color = yellow
	case ERROR:
		levelText = "ERROR"
		color = red
	case FATAL:
		levelText = "FATAL"
		color = purple
	}

	// Format the log entry
	var logEntry string
	if l.ShowTime {
		logEntry = fmt.Sprintf("%s%s [%s] %s%s\n", color, timestamp, levelText, formattedMsg, reset)
	} else {
		logEntry = fmt.Sprintf("%s[%s] %s%s\n", color, levelText, formattedMsg, reset)
	}

	// Write to appropriate output
	if level >= ERROR {
		fmt.Fprint(os.Stderr, logEntry)
	} else {
		fmt.Fprint(os.Stdout, logEntry)
	}

	// Exit program if it's a fatal error
	if level == FATAL {
		os.Exit(1)
	}
}

// Convenience methods for each log level
func (l *Logger) Debug(message string, args ...interface{}) {
	l.Log(DEBUG, message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.Log(INFO, message, args...)
}

func (l *Logger) Warning(message string, args ...interface{}) {
	l.Log(WARNING, message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.Log(ERROR, message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.Log(FATAL, message, args...)
}