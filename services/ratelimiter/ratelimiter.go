package ratelimiter

import (
	"context"
	"fmt"
)

// ex:
// 1. allow 60 requests in 60 seconds
//	count = 60, period = 60
// 2. allow 100 requests in 1 seconds
//	count = 100, period = 1
type Validator struct {
	Count  int32
	Period int32
}

type ValidatorName string

const (
	ValidatorNothing       ValidatorName = "nothing"
	ValidatorGlobalIPLimit ValidatorName = "global-IPLimitation"
)

var Validators = map[ValidatorName]Validator{
	ValidatorGlobalIPLimit: Validator{
		Count:  10,
		Period: 10,
	},
}

var (
	ErrValidatorNotExist = fmt.Errorf("validator doesn't exist")
)

type Service interface {
	Validate(ctx context.Context, name ValidatorName, refArg string) (ok bool, current, remaining int32, err error)
}
