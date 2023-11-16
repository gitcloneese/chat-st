package server

import (
	"time"
	"x-server/core/apollo"
	"x-server/core/pkg/util"
	"x-server/logger/internal/service"
	logger "xy3-proto/logger"
	"xy3-proto/pkg/log"
	"xy3-proto/pkg/net/rpc/warden"

	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"github.com/go-kratos/kratos/v2/transport/grpc"

	"xy3-proto/pkg/conf/paladin"
)

// NewGRPCServer
// New new a grpc server.
func NewGRPCServer(svc *service.Service) (srv *grpc.Server, err error) {
	ct := paladin.TOML{}
	cfg := &warden.ServerConfig{}
	v := apollo.Get(apollo.GrpcNS)
	if v == nil || v.Unmarshal(&ct) != nil {
		if err = paladin.Get(apollo.GrpcNS).Unmarshal(&ct); err != nil {
			log.Error("NewGrpcServer failed: %v", err)
			return nil, err
		}
	}
	if err = ct.Get("Server").UnmarshalTOML(cfg); err != nil {
		return nil, err
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			util.GrpcRecovery(),
			util.GrpcRequestLog(),
		),
	}
	if cfg.Network != "" {
		opts = append(opts, grpc.Network(cfg.Network))
	}
	if cfg.Addr != "" {
		opts = append(opts, grpc.Address(cfg.Addr))
	}
	if cfg.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(cfg.Timeout)))
	}
	srv = grpc.NewServer(opts...)
	logger.RegisterLoggerServer(srv, svc)
	return srv, err
}
