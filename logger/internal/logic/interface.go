package logic

import (
	"context"
	pb "xy3-proto/logger"
)

type Handler interface {
	Login(context.Context, *pb.LogMsg) error
	Register(context.Context, *pb.LogMsg) error
	Logout(context.Context, *pb.LogMsg) error
}
