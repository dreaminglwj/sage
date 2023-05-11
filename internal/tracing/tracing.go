package tracing

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

// TraceID returns a traceid valuer.
func TraceID() log.Valuer {
	return func(ctx context.Context) interface{} {
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			return ""
		}
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			return sc.TraceID()
		}
		return ""
	}
}

// SpanID returns a spanID valuer.
func SpanID() log.Valuer {
	return func(ctx context.Context) interface{} {
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			return ""
		}
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			return sc.SpanID()
		}
		return ""
	}
}
