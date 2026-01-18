// Package telemetry provides OpenTelemetry setup for OmniRoute services.
// It configures tracing, metrics, and logging exporters.
package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds telemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	EnableConsole  bool
	SampleRate     float64
}

// Provider wraps the OpenTelemetry TracerProvider
type Provider struct {
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

// NewProvider creates a new telemetry provider
func NewProvider(ctx context.Context, cfg *Config) (*Provider, error) {
	if cfg.SampleRate == 0 {
		cfg.SampleRate = 1.0 // Default: sample everything
	}

	// Create resource with service info
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
		resource.WithOS(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	// Create exporters
	var exporters []sdktrace.SpanExporter

	// OTLP exporter (if configured)
	if cfg.OTLPEndpoint != "" {
		otlpExporter, err := otlptrace.New(ctx,
			otlptracegrpc.NewClient(
				otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
				otlptracegrpc.WithInsecure(),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("create OTLP exporter: %w", err)
		}
		exporters = append(exporters, otlpExporter)
	}

	// Console exporter (for development)
	if cfg.EnableConsole {
		consoleExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("create console exporter: %w", err)
		}
		exporters = append(exporters, consoleExporter)
	}

	// If no exporters, use noop
	if len(exporters) == 0 {
		return &Provider{
			tracer: otel.Tracer(cfg.ServiceName),
		}, nil
	}

	// Create span processors
	var opts []sdktrace.TracerProviderOption
	opts = append(opts, sdktrace.WithResource(res))
	opts = append(opts, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)))

	for _, exp := range exporters {
		opts = append(opts, sdktrace.WithBatcher(exp,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
		))
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(opts...)

	// Set global provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{
		tracerProvider: tp,
		tracer:         tp.Tracer(cfg.ServiceName),
	}, nil
}

// Tracer returns the tracer for the service
func (p *Provider) Tracer() trace.Tracer {
	return p.tracer
}

// Shutdown gracefully shuts down the telemetry provider
func (p *Provider) Shutdown(ctx context.Context) error {
	if p.tracerProvider != nil {
		return p.tracerProvider.Shutdown(ctx)
	}
	return nil
}

// StartSpan starts a new span with the given name
func (p *Provider) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return p.tracer.Start(ctx, name, opts...)
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetError records an error on the current span
func SetError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// Common attribute helpers
func ServiceAttr(name string) attribute.KeyValue {
	return attribute.String("service.name", name)
}

func TenantAttr(tenantID string) attribute.KeyValue {
	return attribute.String("tenant.id", tenantID)
}

func UserAttr(userID string) attribute.KeyValue {
	return attribute.String("user.id", userID)
}

func RequestIDAttr(requestID string) attribute.KeyValue {
	return attribute.String("request.id", requestID)
}

func HTTPMethodAttr(method string) attribute.KeyValue {
	return attribute.String("http.method", method)
}

func HTTPRouteAttr(route string) attribute.KeyValue {
	return attribute.String("http.route", route)
}

func HTTPStatusAttr(code int) attribute.KeyValue {
	return attribute.Int("http.status_code", code)
}

func DBSystemAttr(system string) attribute.KeyValue {
	return attribute.String("db.system", system)
}

func DBStatementAttr(statement string) attribute.KeyValue {
	return attribute.String("db.statement", statement)
}

func MessagingSystemAttr(system string) attribute.KeyValue {
	return attribute.String("messaging.system", system)
}

func MessagingDestinationAttr(destination string) attribute.KeyValue {
	return attribute.String("messaging.destination", destination)
}
