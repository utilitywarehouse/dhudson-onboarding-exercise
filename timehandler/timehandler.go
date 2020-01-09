package timehandler

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger, counter prometheus.Counter, prefix, timeFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		counter.Inc()

		tf := time.Now().Format(timeFormat)
		logger.WithField("time_resp", tf).Debug("Time requested")
		_, err := w.Write([]byte(prefix + tf))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("failed to write response"))
			if err != nil {
				logger.WithError(err).Warn("Failed to report failed response")
			}
		}
	}
}
