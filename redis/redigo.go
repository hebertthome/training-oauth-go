package redis

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"

	"github.com/garyburd/redigo/redis"
)

type rdgo struct {
	LcCache
	conf config.RedigoConfig
	log  logger.Logger
}

func startRedigo(c *config.RedigoConfig) LcCache {
	log := logger.GetLogger()
	return cast(&rdgo{conf: *c, log: log})
}

func (rd rdgo) redisConnect() redis.Conn {
	c, err := redis.Dial("tcp", rd.conf.Address)
	if err != nil {
		rd.log.Debug("Redigo",
			logger.String("Address", rd.conf.Address),
		)
		rd.log.Panic("We have a problem with the Redis Connection, plz check the configuration!")
	}
	return c
}

func (rd *rdgo) Del(key string) error {
	conn := rd.redisConnect()
	defer conn.Close()

	if _, err := redis.Bytes(conn.Do("DEL", key)); err != nil {
		rd.log.Error("Redigo",
			logger.String("Del", fmt.Sprintf("Problem on delete key <%v>", key)),
			logger.Struct("Error", err),
		)
		return fmt.Errorf("Problem on delete key <%v>", key)
	}

	return nil
}

func (rd *rdgo) Set(key string, value LcCached) error {
	conn := rd.redisConnect()
	defer conn.Close()

	valueB, err := json.Marshal(value)
	if err != nil {
		rd.log.Error("Redigo",
			logger.String("Set", fmt.Sprintf("Problem on marshal object")),
			logger.Struct("Error", err),
		)
		return fmt.Errorf("Problem on marshal object")
	}

	ttl, err := conn.Do("TTL", key)
	if err != nil || ttl.(int64) < 0 {
		ttl = rd.conf.TtlInSeconds
	}

	if _, err := conn.Do("SETEX", key, ttl, valueB); err != nil {
		rd.log.Error("Redigo",
			logger.String("Set", fmt.Sprintf("Problem on write object on cache by key <%v>", key)),
			logger.Struct("Error", err),
		)
		return fmt.Errorf("Problem on write object on cache by key <%v>", key)
	}

	return err
}

func (rd *rdgo) Get(key string) (*LcCached, error) {
	conn := rd.redisConnect()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		rd.log.Error("Redigo",
			logger.String("Get", fmt.Sprintf("Problem on read object from cache by key <%v>", key)),
			logger.Struct("Error", err),
		)
		return nil, fmt.Errorf("Problem on read object from cache by key <%v>", key)
	}

	var c LcCached
	if err := json.Unmarshal(data, &c); err != nil {
		rd.log.Error("Redigo",
			logger.String("Get", fmt.Sprintf("Problem on unmarshal object")),
			logger.Struct("Error", err),
		)
		return nil, fmt.Errorf("Problem on unmarshal object")
	}

	if c.Count >= rd.conf.ExpireCount {
		rd.log.Error("Redigo",
			logger.String("Get", fmt.Sprintf("Exipre token by number of use: %v/%v", c.Count, rd.conf.ExpireCount)),
			logger.Struct("Error", err),
		)
		return nil, fmt.Errorf("Expired token")
	}

	return &c, nil
}
