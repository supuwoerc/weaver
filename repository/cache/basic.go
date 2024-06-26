package cache

import (
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"github.com/gin-gonic/gin"
)

type BasicCache struct {
	redis *redis.RedisClient
	ctx   *gin.Context
}

var basicCache *BasicCache

func NewBasicCache(ctx *gin.Context) *BasicCache {
	if basicCache == nil {
		basicCache = &BasicCache{
			redis: global.RedisClient,
			ctx:   ctx,
		}
	}
	return basicCache
}
