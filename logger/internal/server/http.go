package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"
	"x-server/core/apollo"
	"x-server/core/pkg/errorEncoder"
	"x-server/core/pkg/util"
	"x-server/logger/internal/service"
	logger "xy3-proto/logger"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
	cfg "xy3-proto/pkg/net/http/config"
)

// NewHTTPServer New new a bm server.
func NewHTTPServer(svc *service.Service) (srv *http.Server, err error) {
	cfgs := &cfg.ServerConfig{}
	ct := paladin.TOML{}
	v := apollo.Get(apollo.HttpNS)
	if v == nil || v.Unmarshal(&ct) != nil {
		if err = paladin.Get(apollo.HttpNS).Unmarshal(&ct); err != nil {
			log.Error("NewHTTPServer failed: %v", err)
			return nil, err
		}
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfgs); err != nil {
		return nil, err
	}
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			util.Recovery(),
		),
	}
	if cfgs.Network != "" {
		opts = append(opts, http.Network(cfgs.Network))
	}
	if cfgs.Addr != "" {
		opts = append(opts, http.Address(cfgs.Addr))
	}
	if cfgs.Timeout != 0 {
		opts = append(opts, http.Timeout(time.Duration(cfgs.Timeout)))
	}
	opts = append(opts, http.ErrorEncoder(errorEncoder.DefaultErrorEncoder))

	srv = http.NewServer(opts...)
	logger.RegisterLoggerHTTPServer(srv, svc)
	return srv, nil
}
