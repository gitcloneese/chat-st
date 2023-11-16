package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"xy3-proto/gomsg/pkg/util"
	"xy3-proto/pkg/log"
	"xy3-proto/trace"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

type playerKey struct{}
type pathKey struct{}
type ipKey struct{}

func Recovery() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				httptr, ok := tr.(http.Transporter)
				if ok {
					params := httptr.Request().URL.Query()
					id := params.Get("userid")
					playerId, _ := strconv.ParseInt(id, 10, 64)
					ctx = context.WithValue(ctx, playerKey{}, playerId)
					ctx = context.WithValue(ctx, pathKey{}, httptr.Request().URL.Path)
					ip := tr.RequestHeader().Get("x-forwarded-for")
					ips := strings.Split(ip, ",")
					ctx = context.WithValue(ctx, ipKey{}, ips[0])
				}
			}
			return handler(ctx, req)
		}
	}
}

func FromContextPlayerId(ctx context.Context) int64 {
	playerId := ctx.Value(playerKey{})
	if playerId != nil {
		return playerId.(int64)
	}
	return 0
}

func FromContextPath(ctx context.Context) string {
	path := ctx.Value(pathKey{})
	if path != nil {
		return path.(string)
	}
	return ""
}

func FromContextIp(ctx context.Context) string {
	ip := ctx.Value(ipKey{})
	if ip != nil {
		return ip.(string)
	}
	return ""
}

var encoder encoding.Codec

func init() {
	encoder = encoding.GetCodec("json")
	if encoder == nil {
		encoder = encoding.GetCodec(proto.Name)
	}
}

// HttpRequestLog
// TODO 待优化
func HttpRequestLog() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			caller := FromContextPlayerId(ctx)
			traceId := caller
			if traceId == 0 {
				traceId = startTime.UnixNano()
			}
			ctx = context.WithValue(ctx, trace.TraceKey{}, traceId)

			reply, err = handler(ctx, req)

			path := FromContextPath(ctx)
			duration := time.Since(startTime)
			logFields := []log.D{
				log.KVInt64(log.KeyUser, caller),
				log.KVString(log.KeyPath, path),
				log.KV("error", err),
				//log.KVInt64(log.TraceId, traceId),
				log.KVFloat64(log.KeyTS, util.Round(duration.Seconds(), 4)),
				log.KVString("source", "http-access-log"),
			}
			reqLen := 0
			if req != nil {
				v, ok := req.(fmt.Stringer)
				if ok {
					strReq := v.String()
					req1, _ := encoder.Marshal(req)
					reqLen = len(strReq)
					// if reqLen > log.MaxLogMsgLength {
					// 	strReq = strReq[:log.MaxLogMsgLength]
					// }
					logFields = append(logFields, log.KV("args", string(req1)))
				}
			}
			logFields = append(logFields, log.KVInt(log.KeyReqLen, reqLen))

			if reply != nil {
				rep, _ := encoder.Marshal(reply)
				replyLen := len(rep)
				if replyLen > log.MaxLogMsgLength {
					rep = rep[:log.MaxLogMsgLength]
				}
				logFields = append(logFields, log.KV(log.KeyRet, string(rep)))
				logFields = append(logFields, log.KVInt(log.KeyRspLen, replyLen))
			}

			logFn(err, duration)(ctx, logFields...)
			return reply, err
		}
	}
}

func GrpcRecovery() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				grpcptr, ok := tr.(*grpc.Transport)
				if ok {
					ctx = context.WithValue(ctx, pathKey{}, grpcptr.Operation())
				}
			}
			return handler(ctx, req)
		}
	}
}

// GrpcRequestLog
// TODO 待优化
func GrpcRequestLog() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			traceId := startTime.UnixNano()
			ctx = context.WithValue(ctx, trace.TraceKey{}, traceId)

			reply, err = handler(ctx, req)

			path := FromContextPath(ctx)
			duration := time.Since(startTime)
			logFields := []log.D{
				log.KVString(log.KeyPath, path),
				log.KV("error", err),
				//log.KVInt64(log.TraceId, traceId),
				log.KVFloat64(log.KeyTS, util.Round(duration.Seconds(), 4)),
				log.KVString("source", "grpc-access-log"),
			}
			reqLen := 0
			if req != nil {
				v, ok := req.(fmt.Stringer)
				if ok {
					strReq := v.String()
					req1, _ := encoder.Marshal(req)
					reqLen = len(strReq)
					// if reqLen > log.MaxLogMsgLength {
					// 	strReq = strReq[:log.MaxLogMsgLength]
					// }
					logFields = append(logFields, log.KV("args", string(req1)))
				}
			}
			logFields = append(logFields, log.KVInt(log.KeyReqLen, reqLen))

			if reply != nil {
				rep, _ := encoder.Marshal(reply)
				replyLen := len(rep)
				if replyLen > log.MaxLogMsgLength {
					rep = rep[:log.MaxLogMsgLength]
				}
				logFields = append(logFields, log.KV(log.KeyRet, string(rep)))
				logFields = append(logFields, log.KVInt(log.KeyRspLen, replyLen))
			}

			logFn(err, duration)(ctx, logFields...)
			return reply, err
		}
	}
}

func logFn(err error, dt time.Duration) func(context.Context, ...log.D) {
	switch {
	case err != nil:
		return log.Errorv
	case dt >= time.Millisecond*500:
		return log.Warnv
	}
	return log.Infov
}
