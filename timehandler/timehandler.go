package timehandler

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger, prefix, timeFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		tf := time.Now().Format(timeFormat)
		logger.WithField("time_resp", tf).Debug("Time requested")
		w.Write([]byte(prefix + tf))
	}
}
