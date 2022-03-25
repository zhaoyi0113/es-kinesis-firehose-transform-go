package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyi0113/es-kinesis-firehose-transform-go/internal"
)

func CreateRoute() *gin.Engine {
	var r = gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		message := internal.Ping()
		c.JSON(200, gin.H{
			"message": message,
		})
	})
	r.POST("/logs", func(c *gin.Context) {

	})
	return r
}
