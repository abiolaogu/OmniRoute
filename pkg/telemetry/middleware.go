// Package telemetry provides Gin middleware for OpenTelemetry tracing.
package telemetry

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// GinMiddleware returns a Gin middleware for tracing HTTP requests
func GinMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Extract trace context from incoming request
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		if c.FullPath() == "" {
			spanName = fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request.Method),
				semconv.HTTPTargetKey.String(c.Request.URL.Path),
				semconv.HTTPURLKey.String(c.Request.URL.String()),
				semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
				semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
				semconv.NetHostNameKey.String(c.Request.Host),
			),
		)
		defer span.End()

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Record response attributes
		status := c.Writer.Status()
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(status),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
			attribute.Float64("http.duration_ms", float64(time.Since(start).Milliseconds())),
		)

		// Record errors
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("error.message", c.Errors.String()))
			for _, err := range c.Errors {
				span.RecordError(err.Err)
			}
		}

		// Set span status based on HTTP status
		if status >= 500 {
			span.SetStatus(1, "Internal Server Error") // Error
		} else if status >= 400 {
			span.SetStatus(1, "Client Error") // Error
		}
	}
}

// TracingContextMiddleware adds trace ID to response headers
func TracingContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().HasTraceID() {
			c.Header("X-Trace-ID", span.SpanContext().TraceID().String())
		}
		c.Next()
	}
}

// RequestIDMiddleware adds a request ID to the context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Add to span
		span := trace.SpanFromContext(c.Request.Context())
		span.SetAttributes(attribute.String("request.id", requestID))

		c.Next()
	}
}
