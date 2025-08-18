package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/supuwoerc/weaver/conf"
	local "github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	storeRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

type LimiterMiddleware struct {
	client redis.UniversalClient
	conf   *conf.Config
}

func NewLimiterMiddleware(redisClient *local.CommonRedisClient, conf *conf.Config) *LimiterMiddleware {
	return &LimiterMiddleware{
		client: redisClient.Client,
		conf:   conf,
	}
}

func (r *LimiterMiddleware) RequestLimit(args ...string) gin.HandlerFunc {
	pattern := r.conf.System.RateLimit.Pattern
	if len(args) > 0 && args[0] != "" {
		pattern = args[0]
	}
	storePrefix := r.conf.System.RateLimit.Prefix
	if len(args) > 1 && args[1] != "" {
		storePrefix = args[1]
	}
	rate, err := limiter.NewRateFromFormatted(pattern)
	if err != nil {
		panic(errors.WithMessage(err, "rate-limiter pattern invalid error"))
	}
	store, err := storeRedis.NewStoreWithOptions(r.client, limiter.StoreOptions{
		Prefix: storePrefix,
	})
	if err != nil {
		panic(errors.WithMessage(err, "rate-limiter store init error"))
	}
	return ginLimiter.NewMiddleware(limiter.New(store, rate), ginLimiter.WithErrorHandler(func(c *gin.Context, err error) {
		response.FailWithError(c, err)
	}), ginLimiter.WithLimitReachedHandler(func(c *gin.Context) {
		response.FailWithError(c, response.Busy)
	}))
}
