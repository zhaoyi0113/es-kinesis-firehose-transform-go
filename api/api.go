package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zhaoyi0113/es-kinesis-firehose-transform-go/internal"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func CreateRoute() *gin.Engine {
	var r = gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		message := internal.Ping()
		c.JSON(200, gin.H{
			"message": message,
		})
	})

	r.POST("/logs", func(c *gin.Context) {
		fmt.Println("v1 receive log event")
		PrintMemUsage()
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
