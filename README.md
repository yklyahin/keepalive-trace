# KeepAlive connection tracing

Writes statistic of keep-alive connections for prometheus.


```go
vec := prometheus.NewCounterVec(prometheus.CounterOpts{
    Name: "http_keepalive",
}, []string{"service", "reused"})


client := http.Client{
 Transport: keepalivetrace.WithRoundTripper(
        http.DefaultTransport,
        keepalivetrace.NewPrometheusTracer("test.service", rate, vec),
 ),
}

client.Do(...)
```