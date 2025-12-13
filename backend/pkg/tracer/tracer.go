package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Setup はTracerProviderを構築し、OpenTelemetryのグローバル設定を更新します。
func Setup(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	cfg.normalize()

	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(buildResourceAttributes(cfg)...))
	if err != nil {
		return nil, nil, fmt.Errorf("build resource: %w", err)
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRatio))
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	cleanup := func(stopCtx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(stopCtx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(shutdownCtx)
	}

	return tp, cleanup, nil
}

func buildResourceAttributes(cfg Config) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, semconv.DeploymentEnvironment(cfg.Environment))
	}
	for k, v := range cfg.ResourceAttributes {
		attrs = append(attrs, attribute.String(k, v))
	}
	return attrs
}

func newExporter(ctx context.Context, cfg Config) (sdktrace.SpanExporter, error) {
	switch cfg.Protocol {
	case ProtocolHTTP:
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
			otlptracehttp.WithTimeout(cfg.Timeout),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		if len(cfg.Headers) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(cfg.Headers))
		}
		return otlptracehttp.New(ctx, opts...)
	case ProtocolGRPC:
		fallthrough
	default:
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
			otlptracegrpc.WithTimeout(cfg.Timeout),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		if len(cfg.Headers) > 0 {
			opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
		}
		return otlptracegrpc.New(ctx, opts...)
	}
}

// Tracer は名前付きトレーサを返します。
func Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name, opts...)
}
