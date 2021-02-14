package logger

import (
	"log"
	"os"
	"strings"
)

type loggingService struct {
	logLevel int
}

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

var Logger loggingService

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	Logger = loggingService{logLevel: WARNING}

	val, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		Logger.SetLogLevel(val)
	}
}

func (e *loggingService) log(logLevel int, args ...string) {
	if logLevel >= e.logLevel {
		switch logLevel {
		case DEBUG, INFO:
			log.Println(args)
		case WARNING:
			log.Print("Warning: ")
			log.Println(args)
		case ERROR:
			log.Panic(args)
		}
	}
}

func (e *loggingService) Debug(args ...string) {
	e.log(DEBUG, args...)
}

func (e *loggingService) Info(args ...string) {
	e.log(INFO, args...)
}

func (e *loggingService) Warning(args ...string) {
	e.log(WARNING, args...)
}

func (e *loggingService) Error(args ...string) {
	e.log(ERROR, args...)
}

func (e *loggingService) SetLogLevel(inputLogLevel string) {
	logLevel := WARNING
	switch strings.ToUpper(inputLogLevel) {
	case "DEBUG":
		logLevel = DEBUG
	case "INFO":
		logLevel = INFO
	case "WARNING":
		logLevel = WARNING
	case "ERROR":
		logLevel = ERROR
	default:
		log.Println("Illegal log level received, defaulting to WARNING")
	}

	e.logLevel = logLevel
}
