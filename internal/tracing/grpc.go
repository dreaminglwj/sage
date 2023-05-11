package tracing

import (
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return grpc_opentracing.UnaryClientInterceptor()
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return grpc_opentracing.StreamClientInterceptor()
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_opentracing.UnaryServerInterceptor()
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_opentracing.StreamServerInterceptor()
}
