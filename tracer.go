package keepalivetrace

import (
	"math/rand"
	"net/http"
	"net/http/httptrace"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Tracer for keep-alive connection
	// SampleRate 0..100 percent of traced requests
	Tracer struct {
		SampleRate int
		trace      *httptrace.ClientTrace
	}
)

// NewPrometheusTracer tracer for prometheus
func NewPrometheusTracer(name string, rate int, vec *prometheus.CounterVec) *Tracer {
	const (
		newconn    = "new"
		reusedconn = "reused"
	)

	hooks := &httptrace.ClientTrace{
		GotConn: func(con httptrace.GotConnInfo) {
			if con.Reused {
				vec.WithLabelValues(name, reusedconn).Inc()
			} else {
				vec.WithLabelValues(name, newconn).Inc()
			}
		},
	}

	return &Tracer{SampleRate: rate, trace: hooks}
}

// WithRequest wrap HTTP request
func (tracer Tracer) WithRequest(r *http.Request) *http.Request {
	if !tracer.IsSampled() {
		return r
	}
	ctx := r.Context()
	ctx = httptrace.WithClientTrace(ctx, tracer.trace)
	return r.WithContext(ctx)
}

// IsSampled checks that request must be traced
func (tracer Tracer) IsSampled() bool {
	if tracer.SampleRate == 100 {
		return true
	} else if tracer.SampleRate == 0 {
		return false
	}

	return 100-tracer.SampleRate < rand.Intn(100)
}
