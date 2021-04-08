package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/guachikuo/dcard-demo/api"
	"github.com/guachikuo/dcard-demo/api/demo"
	"github.com/guachikuo/dcard-demo/services/ratelimiter"
	"github.com/guachikuo/dcard-demo/services/redis"
)

func main() {
	// set the standard logger formatter
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	gin.SetMode(gin.ReleaseMode)

	rdb, err := redis.New()
	if err != nil {
		logrus.Panic(err)
		return
	}

	rt := ratelimiter.New(rdb)

	router := gin.Default()
	apiRouter := router.Group("/api", api.IPRateLimiter(rt))

	demo.NewHandler(apiRouter)

	logrus.Info("start serving https request")
	router.Run()
}
