package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utilitywarehouse/dhudson-onboarding-exercise/timehandler"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}
	logger.Level = logrus.DebugLevel

	mux := http.NewServeMux()

	mux.HandleFunc("/time", timehandler.New(logger, "Current time: ", time.RFC3339))

	logger.Info("Listening...")
	logger.WithError(http.ListenAndServe(":3030", mux)).Info("Finished serving")
}
