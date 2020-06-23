package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

// Log system: Uber Zap used
var (
	log logger
)

type bookstoreLogger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Error(msg string, err error, tags ...zap.Field)
	Info(msg string, tags ...zap.Field)
}

type logger struct {
	log *zap.Logger
}

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// GetLevel returns the log
func GetLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

// getOutput func
func getOutput() string {
	output := strings.TrimSpace(os.Getenv(envLogOutput))
	if output == "" {
		return "stdout"
	}
	return output
}

// GetLogger func
func GetLogger() bookstoreLogger {
	return log
}

// getLevel func returns the level atomic const
func getLevel() zapcore.Level {
	return zap.InfoLevel
}

// Printf is a function to print the logger Info message
func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		l.Info(format)
	} else {
		l.Info(fmt.Sprintf(format, v...))
	}
}

// Print is a function to print the logger Info message
func (l logger) Print(v ...interface{}) {
	l.Info(fmt.Sprintf("%v", v...))
}

// Info is a function to build the logger Info message
func (l logger) Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	log.log.Sync()
}

// Error is a function to build the logger Error message
func (l logger) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error: ", err))
	log.log.Info(msg, tags...)
	log.log.Sync()
}
