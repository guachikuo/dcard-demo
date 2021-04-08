package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/guachikuo/dcard-demo/services/ratelimiter"
)

func IPRateLimiter(rt ratelimiter.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		ratelimiterName := ratelimiter.ValidatorGlobalIPLimit
		ok, curCnt, remainingCnt, err := rt.Validate(c.Request.Context(), ratelimiterName, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Error",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Current", strconv.FormatInt(int64(curCnt), 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(remainingCnt), 10))
		c.Next()
	}
}
