package log

import (
	"fmt"
	"log"
	"os"
)

const (
	PanicLevel int = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type LogConfig struct {
	logLevel 		int
}

var (
	logger 		*log.Logger
	logConfig 	LogConfig
)

func init() {
	logger = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds)
}

func SetLogLevel(level int) {
	logConfig.logLevel = level
}

func Debugf(format string, args ...interface{}) {
	if logConfig.logLevel >= DebugLevel {
		logger.SetPrefix("[DEBUG]")
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if logConfig.logLevel >= InfoLevel {
		logger.SetPrefix("[INFO]")
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if logConfig.logLevel >= WarnLevel {
		logger.SetPrefix("[WARN]")
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if logConfig.logLevel >= ErrorLevel {
		logger.SetPrefix("[ERROR]")
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	if logConfig.logLevel >= FatalLevel {
		logger.SetPrefix("[FATAL]")
		logger.Output(2, fmt.Sprintf(format, args...))
	}
}