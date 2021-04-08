package keys

import (
	"strings"
)

func RedisKey(args ...string) string {
	return strings.Join(args, ":")
}

// define Prefix of a group of redis keys
const (
	PfxRateLimiter = "rateLimiter"
)
