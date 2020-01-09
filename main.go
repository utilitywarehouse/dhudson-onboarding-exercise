package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/utilitywarehouse/dhudson-onboarding-exercise/timehandler"
)

var (
	counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests_count",
		Help: "Total number of requests received",
	})
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}
	logger.Level = logrus.DebugLevel

	promMux := http.NewServeMux()
	promMux.Handle("/__/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8081", promMux)
		if err != nil {
			logger.WithError(err).Warn("Failed to serve metrics")
			return
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/time", timehandler.New(logger, counter, "Current time: ", time.RFC3339))

	logger.Info("Listening...")
	logger.WithError(http.ListenAndServe(":3030", mux)).Info("Finished serving")
}
