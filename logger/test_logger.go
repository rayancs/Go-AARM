package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// TestResult represents the outcome of a test
type TestResult int

const (
	TEST_PASS TestResult = iota
	TEST_INFO
	TEST_FAIL
)

// TestLogger extends the basic Logger with test-specific functionality
type TestLogger struct {
	*Logger
	testName string
}

// NewTestLogger creates a new TestLogger with the specified test name
func NewTestLogger(testName string) *TestLogger {
	return &TestLogger{
		Logger:   New(),
		testName: testName,
	}
}

// getCallerInfo returns the file and line number of the caller
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}

	// Extract just the filename, not the full path
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", fileName, line)
}

// LogTest logs a test result with appropriate formatting
func (t *TestLogger) LogTest(result TestResult, message string, args ...interface{}) {
	var resultText, color string

	switch result {
	case TEST_PASS:
		resultText = "PASS"
		color = green
	case TEST_INFO:
		resultText = "INFO"
		color = cyan
	case TEST_FAIL:
		resultText = "FAIL"
		color = red
	}

	// Format the message with additional args if provided
	formattedMsg := message
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(message, args...)
	}

	// Get current timestamp
	timestamp := ""
	if t.ShowTime {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	callerInfo := getCallerInfo()

	// Format the log entry
	var logEntry string
	if t.ShowTime {
		logEntry = fmt.Sprintf("%s%s [%s] [%s] %s (%s)%s\n",
			color, timestamp, resultText, t.testName, formattedMsg, callerInfo, reset)
	} else {
		logEntry = fmt.Sprintf("%s[%s] [%s] %s (%s)%s\n",
			color, resultText, t.testName, formattedMsg, callerInfo, reset)
	}

	// Write to appropriate output
	if result == TEST_FAIL {
		fmt.Fprint(os.Stderr, logEntry)
	} else {
		fmt.Fprint(os.Stdout, logEntry)
	}
}

// Convenience methods for test logging
func (t *TestLogger) Pass(message string, args ...interface{}) {
	t.LogTest(TEST_PASS, message, args...)
}

func (t *TestLogger) Info(message string, args ...interface{}) {
	t.LogTest(TEST_INFO, message, args...)
}

func (t *TestLogger) Fail(message string, args ...interface{}) {
	t.LogTest(TEST_FAIL, message, args...)
}
