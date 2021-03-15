package redis

import (
	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
)

type LcCached struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type LcCache interface {
	Del(key string) error
	Get(key string) (*LcCached, error)
	Set(key string, value LcCached) error
}

func Start(c *config.Configuration) LcCache {
	if !c.Redis.Redigo.IsStructureEmpty() {
		return startRedigo(&c.Redis.Redigo)
	}
	if !c.Redis.GoRedis.IsStructureEmpty() {
		return startGoRedis(&c.Redis.GoRedis)
	}
	panic("Can't be impossible to configure REDIS!")
}

func cast(x interface{}) LcCache {
	i, ok := x.(LcCache)
	if !ok {
		log := logger.GetLogger()
		log.Error("LcCache",
			logger.String("Cast", "implementation Cache to interface failed"),
		)
	}
	return i
}
