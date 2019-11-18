package keepalivetrace_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	keepalivetrace "github.com/lobotomist/keepalive-trace"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRoundTripper100Rate(t *testing.T) {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_keepalive",
	}, []string{"service", "reused"})

	runner(t, 10, 100, vec)

	reused, _ := vec.GetMetricWithLabelValues("test.service", "reused")
	newconn, _ := vec.GetMetricWithLabelValues("test.service", "new")

	assert.Equal(t, float64(1), testutil.ToFloat64(newconn), "established connections")
	assert.Equal(t, float64(9), testutil.ToFloat64(reused), "reused connections")
}

func TestRoundTripperZeroRate(t *testing.T) {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_keepalive",
	}, []string{"service", "reused"})

	runner(t, 10, 0, vec)

	reused, _ := vec.GetMetricWithLabelValues("test.service", "reused")
	newconn, _ := vec.GetMetricWithLabelValues("test.service", "new")

	assert.Equal(t, float64(0), testutil.ToFloat64(newconn), "established connections")
	assert.Equal(t, float64(0), testutil.ToFloat64(reused), "reused connections")
}

func runner(t *testing.T, N int, rate int, vec *prometheus.CounterVec) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	client := http.Client{
		Transport: keepalivetrace.WithRoundTripper(
			http.DefaultTransport,
			keepalivetrace.NewPrometheusTracer("test.service", rate, vec),
		),
	}

	for i := 0; i < N; i++ {
		r, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
		res, err := client.Do(r)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}

	srv.Close()
}
