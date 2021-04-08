package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	rdb "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const (
	redisAddr = "redis:6379"
)

type impl struct {
	rdb *rdb.Client
}

func New() (Service, error) {
	ctx := context.Background()

	rdbSrv := rdb.NewClient(&rdb.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	result := rdbSrv.Ping(ctx)
	if _, err := result.Result(); err != nil {
		logrus.WithField("err", err).Error("rdbSrv.Ping() failed in New")
		return nil, err
	}

	return &impl{
		rdb: rdbSrv,
	}, nil
}

func (im *impl) ScriptDo(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	luaScript := redis.NewScript(script)
	ret, err := luaScript.Run(ctx, im.rdb, keys, args...).Result()
	if err != nil {
		logrus.WithField("err", err).Error("luaScript.Run failed in ScriptDo")
		return nil, err
	}
	return ret, nil
}
