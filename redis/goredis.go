package redis

import (
	"math/rand"
	"time"

	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type goRD struct {
	conf  config.GoRedisConfig
	log   logger.Logger
	cache *cache.Cache
}

func startGoRedis(c *config.GoRedisConfig) LcCache {
	log := logger.GetLogger()
	log.Debug("LcCache",
		logger.Struct("Addrs", c.Addresses),
	)
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: c.Addresses,
	})
	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(24, time.Hour),
	})
	return cast(goRD{conf: *c, cache: mycache, log: log})
}

func (rd *goRD) Set(key string, value LcCached) error {
	return rd.cache.Set(&cache.Item{
		Ctx:   *&rd.conf.Context,
		Key:   key,
		Value: value,
		TTL:   (time.Duration(rand.Intn(rd.conf.TtlInSeconds)) * time.Second),
	})
}

func (rd *goRD) Get(key string) (*LcCached, error) {
	var result LcCached
	if err := rd.cache.Get(*&rd.conf.Context, key, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (rd *goRD) Del(key string) error {
	return rd.cache.Delete(*&rd.conf.Context, key)
}
