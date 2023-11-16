package context

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"xy3-proto/pkg/log"
)

// http请求时，能将 ctx 转换成 gin.Context .
// 0:内部grpc请求 -1:解析失败 -2:url没有userid
func FromContextPlayerId(ctx context.Context) int64 {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return 0
	}

	if id, ok := c.GetQueryArray("userid"); ok && len(id) == 1 {
		uid, err := strconv.ParseInt(id[0], 10, 64)
		if err != nil {
			log.Errorc(ctx, "err :%v", err)
			return -1
		}
		return uid
	}

	return -2
}
