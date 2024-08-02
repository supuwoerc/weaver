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

func NewBasicCache(ctx *gin.Context) *BasicCache {
	return &BasicCache{
		redis: global.RedisClient,
		ctx:   ctx,
	}
}
