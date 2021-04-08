package demo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct{}

// NewHandler ...
func NewHandler(rg *gin.RouterGroup) {
	hd := handler{}

	demoGroup := rg.Group("/demo")
	demoGroup.GET("/data", hd.getDemoData)
}

func (hd *handler) getDemoData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
