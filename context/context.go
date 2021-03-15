package context

import (
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
	"bitbucket.org/hebertthome/traning-oauth-go/redis"
)

type AppContext struct {
	Cache  redis.LcCache
	Logger logger.Logger
}
