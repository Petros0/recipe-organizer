package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
)

// Logger provides structured logging with request context
type Logger struct {
	ctx       openruntimes.Context
	url       string
	userID    string
	startTime time.Time
}

// LogEntry represents a structured log message
type LogEntry struct {
	Timestamp string                 `json:"ts"`
	Level     LogLevel               `json:"level"`
	URL       string                 `json:"url,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Component string                 `json:"component,omitempty"`
	Message   string                 `json:"msg"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// NewLogger creates a logger with request context
func NewLogger(ctx openruntimes.Context, url, userID string) *Logger {
	return &Logger{
		ctx:       ctx,
		url:       url,
		userID:    userID,
		startTime: time.Now(),
	}
}

// log writes a structured log entry
func (l *Logger) log(level LogLevel, component, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		URL:       l.url,
		UserID:    l.userID,
		Component: component,
		Message:   message,
		Fields:    fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to plain text if JSON fails
		l.ctx.Log(fmt.Sprintf("[%s] %s: %s", component, level, message))
		return
	}

	if level == LogLevelError {
		l.ctx.Error(string(jsonBytes))
	} else {
		l.ctx.Log(string(jsonBytes))
	}
}

// Info logs an informational message
func (l *Logger) Info(component, message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LogLevelInfo, component, message, f)
}

// Error logs an error message
func (l *Logger) Error(component, message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LogLevelError, component, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(component, message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LogLevelWarn, component, message, f)
}

// Debug logs a debug message
func (l *Logger) Debug(component, message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LogLevelDebug, component, message, f)
}

// WithDuration logs with elapsed time since logger creation
func (l *Logger) WithDuration(component, message string, fields ...map[string]interface{}) {
	duration := time.Since(l.startTime).Milliseconds()
	f := map[string]interface{}{
		"duration_ms": duration,
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			f[k] = v
		}
	}
	l.log(LogLevelInfo, component, message, f)
}
