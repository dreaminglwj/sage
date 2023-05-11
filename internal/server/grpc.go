package server

import (
	"time"

	"github.com/dreaminglwj/sage/internal/tracing"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	pb "github.com/dreaminglwj/sage/proto/sage"

	"github.com/dreaminglwj/sage/internal/conf"
	"github.com/dreaminglwj/sage/internal/plugin/log"
	"github.com/dreaminglwj/sage/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Config,
	logger *log.Logger,
	schemaService *service.SchemaService,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		grpc.UnaryInterceptor(tracing.UnaryServerInterceptor()),
	}
	config := c.Server.GRPC
	if config.Network != "" {
		opts = append(opts, grpc.Network(config.Network))
	}
	if config.Addr != "" {
		opts = append(opts, grpc.Address(config.Addr))
	}
	if config.Timeout > 0 {
		opts = append(opts, grpc.Timeout(time.Duration(config.Timeout)*time.Millisecond))
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterSchemaServer(srv, schemaService)

	return srv
}
