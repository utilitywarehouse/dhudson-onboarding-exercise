package main

import (
	"net/http"
	"time"

	op "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/utilitywarehouse/dhudson-onboarding-exercise/timehandler"
	"go.opencensus.io/stats/view"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}
	logger.Level = logrus.DebugLevel

	exporter, err := op.NewExporter(op.Options{})
	if err != nil {
		logger.WithError(err).Warn("Failed to set up new exporter")
		return
	}

	view.RegisterExporter(exporter)

	err = view.Register(timehandler.DefaultViews...)
	if err != nil {
		logger.WithError(err).Warn("Failed to set up new exporter")
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/time", timehandler.New(logger, "Current time: ", time.RFC3339))

	mux.Handle("/metrics", exporter)

	logger.Info("Listening...")
	logger.WithError(http.ListenAndServe(":3030", mux)).Info("Finished serving")
}
