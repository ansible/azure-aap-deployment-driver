package config

import (
	"io"
	"os"

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
	logFile, err := os.OpenFile(GetEnvironment().LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		bothPlaces := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(bothPlaces)
	} else {
		log.Info("Error opening logfile, will log only to stdout")
	}
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(UTCTextFormatter{formatter})
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
}
