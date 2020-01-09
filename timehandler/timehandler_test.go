package timehandler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utilitywarehouse/dhudson-onboarding-exercise/timehandler"
)

var (
	testCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests_count",
		Help: "Total number of requests received",
	})
)

func TestNew(t *testing.T) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}
	logger.Level = logrus.DebugLevel

	tests := map[string]struct {
		prefix     string
		timeFormat string
	}{
		"emptyPrefixEmptyFormat": {},
		"emptyPrefixRFCFormat": {
			timeFormat: time.RFC3339,
		},
		"prefixEmptyFormat": {
			prefix: "a prefix",
		},
		"prefixFormat": {
			prefix:     "a prefix",
			timeFormat: time.RFC3339,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			th := timehandler.New(logger, testCounter, test.prefix, test.timeFormat)

			recorder := httptest.NewRecorder()

			th.ServeHTTP(recorder, nil)

			resp := recorder.Result()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			rc := resp.Body

			body, err := ioutil.ReadAll(rc)
			require.NoError(t, err)

			err = rc.Close()
			require.NoError(t, err)

			bodyStr := string(body)

			assert.Contains(t, bodyStr, test.prefix)

			timeStr := strings.TrimPrefix(bodyStr, test.prefix)

			_, err = time.Parse(test.timeFormat, timeStr)
			assert.NoError(t, err)
		})
	}
}
