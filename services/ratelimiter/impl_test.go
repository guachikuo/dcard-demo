package ratelimiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/guachikuo/dcard-demo/models/keys"
	mredis "github.com/guachikuo/dcard-demo/services/redis/mocks"
)

var (
	mockCTX = context.Background()
)

type ratelimiterSuite struct {
	suite.Suite

	ratelimiter *impl
	mockRedis   *mredis.Service
}

func (r *ratelimiterSuite) SetupSuite() {
	r.mockRedis = &mredis.Service{}
}

func (r *ratelimiterSuite) TearDownSuite() {
}

func (r *ratelimiterSuite) SetupTest() {
	r.ratelimiter = New(r.mockRedis).(*impl)
}

func (r *ratelimiterSuite) TearDownTest() {
	r.mockRedis.AssertExpectations(r.T())
}

func TestRatelimiterSuite(t *testing.T) {
	suite.Run(t, new(ratelimiterSuite))
}

func (r *ratelimiterSuite) TestValidateNotExist() {
	_, _, _, err := r.ratelimiter.Validate(mockCTX, "test", "")
	r.Require().EqualError(err, ErrValidatorNotExist.Error())
}

func (r *ratelimiterSuite) TestValidateWiredCount() {
	Validators["test"] = Validator{
		Count:  0,
		Period: 60,
	}
	ok, cur, remain, err := r.ratelimiter.Validate(mockCTX, "test", "")
	r.Require().NoError(err)
	r.Require().False(ok)
	r.Require().Equal(int32(0), cur)
	r.Require().Equal(int32(0), remain)
}

func (r *ratelimiterSuite) TestValidateWiredPeriod() {
	Validators["test"] = Validator{
		Count:  1,
		Period: 0,
	}
	ok, cur, remain, err := r.ratelimiter.Validate(mockCTX, "test", "")
	r.Require().NoError(err)
	r.Require().True(ok)
	r.Require().Equal(int32(1), cur)
	r.Require().Equal(int32(0), remain)
}

func (r *ratelimiterSuite) TestValidateGlobalIPLimitPass() {
	key := keys.RedisKey(keys.PfxRatelimiter, string(ValidatorGlobalIPLimit), "localhost")
	r.mockRedis.On("ScriptDo", mockCTX, script, []string{key}, int32(60), int32(60)).Return(
		[]interface{}{1, int64(10)}, nil,
	).Once()

	ok, cur, remain, err := r.ratelimiter.Validate(mockCTX, ValidatorGlobalIPLimit, "localhost")
	r.Require().NoError(err)
	r.Require().True(ok)
	r.Require().Equal(int32(10), cur)
	r.Require().Equal(int32(50), remain)
}

func (r *ratelimiterSuite) TestValidateGlobalIPLimitTooManyRequests() {
	key := keys.RedisKey(keys.PfxRatelimiter, string(ValidatorGlobalIPLimit), "localhost")
	r.mockRedis.On("ScriptDo", mockCTX, script, []string{key}, int32(60), int32(60)).Return(
		[]interface{}{nil, 60}, nil,
	).Once()

	ok, cur, remain, err := r.ratelimiter.Validate(mockCTX, ValidatorGlobalIPLimit, "localhost")
	r.Require().NoError(err)
	r.Require().False(ok)
	r.Require().Equal(int32(60), cur)
	r.Require().Equal(int32(0), remain)
}
