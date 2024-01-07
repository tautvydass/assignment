package log

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var now = time.Now

// Trace logs a trace message.
func Trace(message string) {
	message = formatMessage(message)
	color.White(message)
}

// Tracef logs a formatted trace message.
func Tracef(format string, args ...interface{}) {
	message := formatMessage(fmt.Sprintf(format, args...))
	color.White(message)
}

// Info logs an info message.
func Info(message string) {
	message = formatMessage(message)
	color.Green(message)
}

// Infof logs a formatted info message.
func Infof(format string, args ...interface{}) {
	message := formatMessage(fmt.Sprintf(format, args...))
	color.Green(message)
}

// Warn logs a warning message.
func Warn(message string) {
	message = formatMessage(message)
	color.Yellow(message)
}

// Warnf logs a formatted warning message.
func Warnf(format string, args ...interface{}) {
	message := formatMessage(fmt.Sprintf(format, args...))
	color.Yellow(message)
}

// Error logs an error message.
func Error(message string) {
	message = formatMessage(message)
	color.Red(message)
}

// Errorf logs a formatted error message.
func Errorf(format string, args ...interface{}) {
	message := formatMessage(fmt.Sprintf(format, args...))
	color.Red(message)
}

func formatMessage(message string) string {
	return fmt.Sprintf("[%s] %s\n", now().Format("2006-01-02 15:04:05"), message)
}
