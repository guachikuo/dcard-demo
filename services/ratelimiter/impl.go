package ratelimiter

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/guachikuo/dcard-demo/models/keys"
	"github.com/guachikuo/dcard-demo/services/redis"
)

const (
	script = `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])

		local cnt = redis.call('GET', key)
		if (cnt == false) then
			redis.call('SETEX', key, ttl, 1)
			return {true, 1}
		elseif (tonumber(cnt) == limit) then
			return {false, limit}
		else
			local newCnt = redis.call('INCR', key)
			return {true, newCnt}
		end
	`
)

type impl struct {
	rdb redis.Service
}

func New(
	rdb redis.Service,
) Service {
	return &impl{
		rdb: rdb,
	}
}

func (im *impl) Validate(ctx context.Context, name ValidatorName, refArg string) (ok bool, current, remaining int32, err error) {
	validator, ok := Validators[name]
	if !ok {
		return false, 0, 0, ErrValidatorNotExist
	}

	limitCnt := validator.Count
	if limitCnt == 0 {
		return false, 0, 0, nil
	}

	period := validator.Period
	if period == 0 {
		return true, 1, 0, nil
	}

	key := keys.RedisKey(keys.PfxRatelimiter, string(name), refArg)
	ret, err := im.rdb.ScriptDo(ctx, script, []string{key}, limitCnt, period)
	if err != nil {
		logrus.WithField("err", err).Error("rdb.ScriptDo failed in Validate")
		return false, 0, 0, err
	}

	values := ret.([]interface{})
	ok = values[0] != nil
	if !ok {
		return false, limitCnt, 0, nil
	}

	newCnt := int32(values[1].(int64))
	return true, newCnt, limitCnt - newCnt, nil
}
