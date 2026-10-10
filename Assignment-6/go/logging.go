// The log line format: HH:MM:SS.mmm LEVEL WHO        | message

package main

import (
	"fmt"
	"io"
	"log"
	"time"
)

// Logger wraps the standard log package, which is safe for use by many goroutines at once.
// The log flags are 0, so the program formats the time itself and the line matches the Java format.
type Logger struct {
	out *log.Logger
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{out: log.New(w, "", 0)}
}

func (l *Logger) Info(who, format string, args ...any)  { l.write("INFO", who, format, args) }
func (l *Logger) Warn(who, format string, args ...any)  { l.write("WARN", who, format, args) }
func (l *Logger) Error(who, format string, args ...any) { l.write("ERROR", who, format, args) }

func (l *Logger) write(level, who, format string, args []any) {
	l.out.Printf("%s %-5s %-10s | %s", time.Now().Format("15:04:05.000"), level, who, fmt.Sprintf(format, args...))
}
