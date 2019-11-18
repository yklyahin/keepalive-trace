package keepalivetrace

import (
	"net/http"
)

type transport struct {
	tracer *Tracer
	prev   http.RoundTripper
}

// WithRoundTripper add tracer to the transport
func WithRoundTripper(rt http.RoundTripper, tracer *Tracer) http.RoundTripper {
	return &transport{
		tracer: tracer,
		prev:   rt,
	}
}

func (tr transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return tr.prev.RoundTrip(tr.tracer.WithRequest(r))
}
