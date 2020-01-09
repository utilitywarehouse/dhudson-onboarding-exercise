package timehandler

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var timeCallCountMeasure = stats.Int64("onboarding-dhudson/measure/time_call_count", "Count of times the time has been requested", stats.UnitDimensionless)

// DefaultViews expose the default views for a time handler.
var DefaultViews = []*view.View{
	&view.View{
		Name:        "onboarding-dhudson/views/time_call_count",
		Description: "The total amount of times the time has been requested",
		Measure:     timeCallCountMeasure,
		Aggregation: view.Count(),
	},
}

var timeCallCount int64

func New(logger *logrus.Logger, prefix, timeFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		ctx := context.Background()
		atomic.AddInt64(&timeCallCount, 1)
		stats.Record(ctx, timeCallCountMeasure.M(timeCallCount))

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
