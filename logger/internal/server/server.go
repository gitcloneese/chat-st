package server

import (
	"github.com/google/wire"
	proto "xy3-proto"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, proto.NewRegister)
