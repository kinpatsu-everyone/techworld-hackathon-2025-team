package tracer

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// HTTPHandler はnet/httpハンドラをOpenTelemetryで計測します。
func HTTPHandler(tp trace.TracerProvider, name string, handler http.Handler, opts ...otelhttp.Option) http.Handler {
	if tp == nil {
		tp = otel.GetTracerProvider()
	}
	fullOpts := append(opts, otelhttp.WithTracerProvider(tp))
	return otelhttp.NewHandler(handler, name, fullOpts...)
}

// HTTPClientTransport は外向きHTTPクライアント用のトランスポートを返します。
func HTTPClientTransport(tp trace.TracerProvider, base http.RoundTripper, opts ...otelhttp.Option) http.RoundTripper {
	if tp == nil {
		tp = otel.GetTracerProvider()
	}
	if base == nil {
		base = http.DefaultTransport
	}
	fullOpts := append(opts, otelhttp.WithTracerProvider(tp))
	return otelhttp.NewTransport(base, fullOpts...)
}
