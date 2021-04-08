package redis

import (
	"context"

	rdb "github.com/go-redis/redis/v8"
)

var (
	ErrNotFound = rdb.Nil
)

type Service interface {
	ScriptDo(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}
