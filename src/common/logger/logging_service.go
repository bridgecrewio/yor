package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type loggingService struct {
	logLevel   LogLevel
	stdout     *os.File
	stderr     *os.File
	tempWriter *os.File
	disabled   int32
}

type LogLevel int
type ErrorType int

var MuteLock sync.Mutex

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

const (
	SILENT ErrorType = iota
)

var strLogLevels = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

var strErrorTypes = map[string]ErrorType{
	"SILENT": SILENT,
}

var Logger loggingService

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	Logger = loggingService{logLevel: WARNING, stdout: os.Stdout, stderr: os.Stderr, disabled: 0}

	val, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		Logger.SetLogLevel(val)
	}
}

func (e *loggingService) log(logLevel LogLevel, args ...string) {
	if logLevel >= e.logLevel {
		var strArgs string
		if len(args) == 2 {
			strArgs = strings.Join([]string{args[0]}, " ")

		} else {
			strArgs = strings.Join(args, " ")
		}
		strArgs = fmt.Sprintf("[%s] ", strLogLevels[logLevel]) + strArgs
		switch logLevel {
		case DEBUG, INFO, WARNING:
			if e.disabled == 0 {
				log.Println(strArgs)
			}
		case ERROR:
			if len(args) == 2 {
				errorType := args[1]
				if _, ok := strErrorTypes[errorType]; ok {
					log.Println(strArgs)
				}
			} else {
				log.Println(strArgs)
			}
			os.Exit(1)
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
	if Logger.logLevel >= WARNING {
		if atomic.LoadInt32(&Logger.disabled) == 0 {
			Debug("Mute logging")
			_, Logger.tempWriter, _ = os.Pipe()
			os.Stdout = Logger.tempWriter
			os.Stderr = Logger.tempWriter
			log.SetOutput(Logger.tempWriter)
		}
		atomic.AddInt32(&Logger.disabled, 1)
	}
}

func UnmuteLogging() {
	if Logger.logLevel >= WARNING {
		atomic.AddInt32(&Logger.disabled, -1)
		if atomic.LoadInt32(&Logger.disabled) == 0 {
			if Logger.tempWriter != nil {
				_ = Logger.tempWriter.Close()
			}
			os.Stdout = Logger.stdout
			os.Stderr = Logger.stderr
			log.SetOutput(os.Stderr)
			Debug("Unmute logging")
		}
	}
}
