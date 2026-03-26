package logger

import (
	"fmt"
	"log"
	"os"
	"passkey-server/config"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

var (
	stdoutLog = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	stderrLog = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
)

func color(color string, message string) string {
	return color + message + reset
}

func Debug(message string) {
	if config.LogLevel <= config.LogLevelDebug {
		stdoutLog.Output(1, "["+color(green, "DEBUG")+"]"+message)

	}
}

func Debugf(message string, args ...any) {
	if config.LogLevel <= config.LogLevelDebug {
		stdoutLog.Output(2, fmt.Sprintf("["+color(green, "DEBUG")+"] "+message, args...))
	}
}

func Info(message string) {
	if config.LogLevel <= config.LogLevelInfo {
		stdoutLog.Output(1, "["+color(cyan, "INFO")+" ] "+message)
	}
}

func Infof(message string, args ...any) {
	if config.LogLevel <= config.LogLevelInfo {
		stdoutLog.Output(2, fmt.Sprintf("["+color(cyan, "INFO")+" ] "+message, args...))
	}
}

func Warn(message string) {
	if config.LogLevel <= config.LogLevelWarn {
		stdoutLog.Output(1, "["+color(yellow, "WARN")+" ] "+message)
	}
}

func Warnf(message string, args ...any) {
	if config.LogLevel <= config.LogLevelWarn {
		stdoutLog.Output(2, fmt.Sprintf("["+color(yellow, "WARN")+" ] "+message, args...))
	}
}

func Error(message string) {
	if config.LogLevel <= config.LogLevelError {
		stderrLog.Output(1, "["+color(red, "ERROR")+"] "+message)
	}
}

func Errorf(message string, args ...any) {
	if config.LogLevel <= config.LogLevelError {
		stderrLog.Output(2, fmt.Sprintf("["+color(red, "ERROR")+"] "+message, args...))
	}
}

func Fatal(message string) {
	if config.LogLevel <= config.LogLevelError {
		stderrLog.Output(1, "["+color(red, "FATAL")+"] "+message)
	}
}

func Fatalf(message string, args ...any) {
	if config.LogLevel <= config.LogLevelError {
		stderrLog.Output(2, fmt.Sprintf("["+color(red, "FATAL")+"] "+message, args...))
	}
}
