package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UTCTextFormatter struct {
	log.Formatter
}

func (u UTCTextFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

func ConfigureLogging() {
	logFile, err := os.OpenFile(GetEnvironment().BASE_PATH+GetEnvironment().LOG_REL_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		bothPlaces := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(bothPlaces)
	} else {
		log.Info("Error opening log file, will log only to stdout")
	}
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	formatter.CallerPrettyfier = CallerFormattingFunc
	log.SetFormatter(UTCTextFormatter{formatter})
	log.SetReportCaller(true)
	log.SetLevel(getLogLevel(GetEnvironment().LOG_LEVEL))
}

func getLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "trace":
		return log.TraceLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	}
	return log.InfoLevel
}

func CallerFormattingFunc(frame *runtime.Frame) (function string, file string) {
	function = frame.Function
	file = fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
	return
}
