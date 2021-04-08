package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/guachikuo/dcard-demo/services/ratelimiter"
	mratelimiter "github.com/guachikuo/dcard-demo/services/ratelimiter/mocks"
)

var (
	mockCTX = context.Background()

	mockHandler = func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
)

type middlewaresSuite struct {
	suite.Suite

	router *gin.Engine

	mockRatelimiter *mratelimiter.Service
}

func (r *middlewaresSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	r.mockRatelimiter = &mratelimiter.Service{}
}

func (r *middlewaresSuite) TearDownSuite() {
}

func (r *middlewaresSuite) SetupTest() {
	r.router = gin.Default()
}

func (r *middlewaresSuite) TearDownTest() {
	r.mockRatelimiter.AssertExpectations(r.T())
}

func TestMiddlewaresSuite(t *testing.T) {
	suite.Run(t, new(middlewaresSuite))
}

func (r *middlewaresSuite) TestIPRateLimiterPass() {
	apiRouter := r.router.Group("/api", IPRateLimiter(r.mockRatelimiter))
	apiRouter.GET("/test", mockHandler)

	req := mockHTTPRequest("GET", "/api/test", nil)
	respRecorder := newRecorder()

	r.mockRatelimiter.On(
		"Validate", mock.Anything, ratelimiter.ValidatorGlobalIPLimit, mockIP,
	).Return(
		true, int32(10), int32(50), nil,
	).Once()
	r.router.ServeHTTP(respRecorder, req)

	// check result
	r.Require().Equal(http.StatusOK, respRecorder.Code)
	r.Require().Equal("10", respRecorder.Header().Get("X-RateLimit-Current"))
	r.Require().Equal("50", respRecorder.Header().Get("X-RateLimit-Remaining"))
}

func (r *middlewaresSuite) TestIPRateLimiterTooManyRequests() {
	apiRouter := r.router.Group("/api", IPRateLimiter(r.mockRatelimiter))
	apiRouter.GET("/test", mockHandler)

	req := mockHTTPRequest("GET", "/api/test", nil)
	respRecorder := newRecorder()

	r.mockRatelimiter.On(
		"Validate", mock.Anything, ratelimiter.ValidatorGlobalIPLimit, mockIP,
	).Return(
		false, int32(60), int32(0), nil,
	).Once()
	r.router.ServeHTTP(respRecorder, req)

	// check result
	r.Require().Equal(http.StatusTooManyRequests, respRecorder.Code)
}
