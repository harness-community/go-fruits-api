package utils

import (
	"io"

	"log"

	"github.com/sirupsen/logrus"
)

//LogSetup sets up the logging for the application
func LogSetup(out io.Writer, level string) *logrus.Logger {
	lvl, err := logrus.ParseLevel(level)

	if err != nil {
		log.Printf("Unable to use the %s level, %#v. Defaulting to warning.", level, err)
		lvl = logrus.WarnLevel
	}

	log := &logrus.Logger{
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:15:10",
		},
		Out: out,
		//ReportCaller: true,
		Level: lvl,
	}

	return log
}
