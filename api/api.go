package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		fmt.Println("receive log event")
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		internal.FailOnError(err, "Failed to parse request body")
		var record internal.LogEvent
		json.Unmarshal(jsonData, &record)
		fmt.Println("receive log event:", len(record.Records))
		response := internal.ProcessLogs(record, "logs")
		c.IndentedJSON(http.StatusOK, response)
	})

	r.POST("/traces", func(c *gin.Context) {
		fmt.Println("Get trace")
		jsonData, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("trace json:", string(jsonData))
	})

	r.POST("/metrics", func(c *gin.Context) {
		// jsonData, err := ioutil.ReadAll(c.Request.Body)
		// internal.FailOnError(err, "Failed to parse request body")
		// var record internal.LogEventRecord
		// json.Unmarshal(jsonData, &record)
		// response := internal.ProcessLogs(record, "metrics")
		// c.IndentedJSON(http.StatusOK, response)
	})
	return r
}
