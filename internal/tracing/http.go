package tracing

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	transportHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var httpTag = opentracing.Tag{Key: string(ext.Component), Value: "http"}

func HTTPServerTracingMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if httpTr, ok := tr.(*transportHttp.Transport); ok {
					r := httpTr.Request()
					parentSpanContext, err := opentracing.GlobalTracer().Extract(
						opentracing.HTTPHeaders,
						opentracing.HTTPHeadersCarrier(r.Header))
					if err == nil || err == opentracing.ErrSpanContextNotFound {
						var opts []opentracing.StartSpanOption
						// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
						opts = append(opts, ext.RPCServerOption(parentSpanContext))
						opts = append(opts, httpTag)
						opts = append(opts, opentracing.Tag{Key: string(ext.HTTPUrl), Value: r.RequestURI})
						opts = append(opts, opentracing.Tag{Key: string(ext.HTTPMethod), Value: r.Method})
						opts = append(opts, opentracing.Tag{Key: "http.host", Value: r.Host})
						opts = append(opts, opentracing.Tag{Key: "http.path", Value: r.URL.Path})
						opts = parseXForwardedHeader(r.Header, opts)
						serverSpan := opentracing.GlobalTracer().StartSpan(
							r.URL.Path,
							opts...,
						)
						ctx = opentracing.ContextWithSpan(ctx, serverSpan)
						defer serverSpan.Finish()
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

func parseXForwardedHeader(header http.Header, opts []opentracing.StartSpanOption) []opentracing.StartSpanOption {
	fwdAddress := header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		opts = append(opts, opentracing.Tag{Key: "http.x-forwarded-for", Value: fwdAddress})
	}
	fwdHost := header.Get("X-Forwarded-Host")
	if fwdHost != "" {
		opts = append(opts, opentracing.Tag{Key: "http.x-forwarded-host", Value: fwdHost})
	}
	fwdPort := header.Get("X-Forwarded-Port")
	if fwdPort != "" {
		opts = append(opts, opentracing.Tag{Key: "http.x-forwarded-port", Value: fwdPort})
	}
	fwdProto := header.Get("X-Forwarded-Proto")
	if fwdProto != "" {
		opts = append(opts, opentracing.Tag{Key: "http.x-forwarded-proto", Value: fwdProto})
	}
	fwdServer := header.Get("X-Forwarded-Server")
	if fwdServer != "" {
		opts = append(opts, opentracing.Tag{Key: "http.x-forwarded-server", Value: fwdServer})
	}
	return opts
}
