package logger

import (
	"log"

	"github.com/tsybulin/gosha/evt"
)

const Topic = "main:logger"

type Severity int8

const (
	LevelDebug Severity = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSystem
)

func (s Severity) String() string {
	return [...]string{"Debug", "Info", "Warn", "Errror", "System"}[s]
}

// Logger ...
type Logger struct {
	severity Severity
}

func (l *Logger) log(v ...interface{}) {
	if severity, ok := v[0].(Severity); ok && severity < l.severity {
		return
	}
	log.Println(v...)
}

func (l *Logger) SetSeverity(severity Severity) {
	l.severity = severity
}

// NewLogger ...
func NewLogger(eventBus evt.Bus) *Logger {
	l := &Logger{
		severity: LevelInfo,
	}
	eventBus.SubscribeAsync(Topic, "Logger.log", l.log, true)
	return l
}
