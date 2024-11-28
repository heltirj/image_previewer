package logger

import (
	"bytes"
	"strings"
	"testing"
)

func setUp() *Logger {
	var buf bytes.Buffer
	logger := New(LogLevelDebug)
	logger.logger.SetOutput(&buf) // Перенаправляем вывод логгера в буфер

	return logger
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevelDebug)
	logger.logger.SetOutput(&buf)

	logger.Debug("This is a debug message")
	output := buf.String()
	if !strings.Contains(output, "[DEBUG] This is a debug message") {
		t.Errorf("Expected debug log not found. Got: %s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevelInfo)
	logger.logger.SetOutput(&buf)

	logger.Info("This is an info message")
	output := buf.String()
	if !strings.Contains(output, "[INFO]  This is an info message") {
		t.Errorf("Expected info log not found. Got: %s", output)
	}
}

func TestLogger_Warning(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevelWarn)
	logger.logger.SetOutput(&buf)

	logger.Warn("This is a warning message")
	output := buf.String()
	if !strings.Contains(output, "[WARN] This is a warning message") {
		t.Errorf("Expected warning log not found. Got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevelError)
	logger.logger.SetOutput(&buf)

	logger.Error("This is an error message")
	output := buf.String()
	if !strings.Contains(output, "[ERROR]  This is an error message") {
		t.Errorf("Expected error log not found. Got: %s", output)
	}
}

func TestLogger_LogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevelWarn)
	logger.logger.SetOutput(&buf)

	logger.Debug("This debug message should not appear")
	logger.Info("This info message should not appear")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	output := buf.String()
	if strings.Contains(output, "[DEBUG]") || strings.Contains(output, "[INFO]") {
		t.Errorf("Unexpected log output found. Got: %s", output)
	}
	if !strings.Contains(output, "[WARN]") {
		t.Errorf("Expected warning log not found. Got: %s", output)
	}
	if !strings.Contains(output, "[ERROR]  This is an error message") {
		t.Errorf("Expected error log not found. Got: %s", output)
	}
}

func TestLogger_DebugKV(t *testing.T) {
	logger := setUp()
	logger.DebugKV("Debugging user login", "user", "john_doe", "attempt", 1)

	output := logger.logger.Writer().(*bytes.Buffer).String()
	expectedOutput := "[DEBUG] Debugging user login [user: john_doe] [attempt: 1]\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected: %q, Got: %q", expectedOutput, output)
	}
}

func TestLogger_InfoKV(t *testing.T) {
	logger := setUp()
	logger.InfoKV("User logged in", "user", "john_doe")

	output := logger.logger.Writer().(*bytes.Buffer).String()
	expectedOutput := "[INFO]  User logged in [user: john_doe]\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected: %q, Got: %q", expectedOutput, output)
	}
}

func TestLogger_WarningKV(t *testing.T) {
	logger := setUp()
	logger.WarnKV("User failed to log in", "user", "john_doe", "reason", "invalid password")

	output := logger.logger.Writer().(*bytes.Buffer).String()
	expectedOutput := "[WARN] User failed to log in [user: john_doe] [reason: invalid password]\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected: %q, Got: %q", expectedOutput, output)
	}
}

func TestLogger_ErrorKV(t *testing.T) {
	logger := setUp()
	logger.ErrorKV("Error while processing request", "user", "john_doe", "error_code", 500)

	output := logger.logger.Writer().(*bytes.Buffer).String()
	expectedOutput := "[ERROR]  Error while processing request [user: john_doe] [error_code: 500]\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected: %q, Got: %q", expectedOutput, output)
	}
}

func TestLogger_IncompleteKeyValue(t *testing.T) {
	logger := setUp()
	logger.ErrorKV("Error occurred", "user", "john_doe", "extra")

	output := logger.logger.Writer().(*bytes.Buffer).String()
	expectedOutput := "[ERROR]  Error occurred [user: john_doe]\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected: %q, Got: %q", expectedOutput, output)
	}
}
