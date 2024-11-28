package logger

import (
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var logLevelStrings = map[string]LogLevel{
	"DEBUG": LogLevelDebug,
	"INFO":  LogLevelInfo,
	"WARN":  LogLevelWarn,
	"ERROR": LogLevelError,
}

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

func New(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Debug(msg string) {
	l.Log(LogLevelDebug, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(LogLevelInfo, msg)
}

func (l *Logger) Warn(msg string) {
	l.Log(LogLevelWarn, msg)
}

func (l *Logger) Error(msg string) {
	l.Log(LogLevelError, msg)
}

func (l *Logger) Log(level LogLevel, message string) {
	if level >= l.level {
		l.logger.Println(formatLogLevel(level) + message)
	}
}

func (l *Logger) DebugKV(msg string, keysAndValues ...interface{}) {
	l.Log(LogLevelDebug, formatMessage(msg, keysAndValues))
}

func (l *Logger) InfoKV(msg string, keysAndValues ...interface{}) {
	l.Log(LogLevelInfo, formatMessage(msg, keysAndValues))
}

func (l *Logger) WarnKV(msg string, keysAndValues ...interface{}) {
	l.Log(LogLevelWarn, formatMessage(msg, keysAndValues))
}

func (l *Logger) ErrorKV(msg string, keysAndValues ...interface{}) {
	l.Log(LogLevelError, formatMessage(msg, keysAndValues))
}

func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var level string
	if err := unmarshal(&level); err != nil {
		return err
	}

	*l = logLevelStrings[level]
	if *l == 0 {
		return fmt.Errorf("invalid log level: %s", level)
	}

	return nil
}

func formatMessage(msg string, keysAndValues []interface{}) string {
	if len(keysAndValues) == 0 {
		return msg
	}

	formatted := msg
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := fmt.Sprintf("%v", keysAndValues[i])
			value := fmt.Sprintf("%v", keysAndValues[i+1])
			formatted += fmt.Sprintf(" [%s: %s]", key, value)
		}
	}
	return formatted
}

func formatLogLevel(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "[DEBUG] "
	case LogLevelInfo:
		return "[INFO]  "
	case LogLevelWarn:
		return "[WARN] "
	case LogLevelError:
		return "[ERROR]  "
	default:
		return ""
	}
}
