package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type loggingService struct {
	logLevel   LogLevel
	stdout     *os.File
	stderr     *os.File
	tempWriter *os.File
	enabled    bool
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

var strLogLevels = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

var Logger loggingService

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	Logger = loggingService{logLevel: WARNING, stdout: os.Stdout, stderr: os.Stderr, enabled: true}

	val, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		Logger.SetLogLevel(val)
	}
}

func (e *loggingService) log(logLevel LogLevel, args ...string) {
	if logLevel >= e.logLevel {
		strArgs := strings.Join(args, " ")
		strArgs = fmt.Sprintf("[%s] ", strLogLevels[logLevel]) + strArgs
		switch logLevel {
		case DEBUG, INFO, WARNING:
			if e.enabled {
				log.Println(strArgs)
			}
		case ERROR:
			log.Panic(strArgs)
		}
	}
}

func Debug(args ...string) {
	Logger.log(DEBUG, args...)
}

func Info(args ...string) {
	Logger.log(INFO, args...)
}

func Warning(args ...string) {
	Logger.log(WARNING, args...)
}

func Error(args ...string) {
	Logger.log(ERROR, args...)
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

func MuteLogging() {
	Debug("Mute logging")
	_, Logger.tempWriter, _ = os.Pipe()
	os.Stdout = Logger.tempWriter
	os.Stderr = Logger.tempWriter
	log.SetOutput(Logger.tempWriter)
	Logger.enabled = false
}

func UnmuteLogging() {
	if Logger.tempWriter != nil {
		Logger.tempWriter.Close()
	}
	os.Stdout = Logger.stdout
	os.Stderr = Logger.stderr
	log.SetOutput(os.Stderr)
	Logger.enabled = true
	Debug("Unmute logging")
}
