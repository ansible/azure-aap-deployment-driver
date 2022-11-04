package config

import (
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
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(UTCTextFormatter{formatter})
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
}
