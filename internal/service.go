package internal

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func Ping() string {
	return "Pong"
}

type LogEventRecord struct {
	Req struct {
		Body LogEvent `json:"body"`
	} `json:"req"`
}

type LogEvent struct {
	Records []struct {
		Data string `json:"data"`
	} `json:"records"`
	Timestamp int    `json:"timestamp"`
	RequestId string `json:"requestId"`
}

func ProcessLogs(event LogEventRecord, t string) {
	for _, record := range event.Req.Body.Records {
		if t == "logs" {
			log := decodeLogEvent(record.Data)
			fmt.Println("xxxx:", len(log))
		}
	}
}

func getIndexName(indexType string) string {
	now := time.Now()
	return fmt.Sprintf("%s-%s", indexType, now.Format("2006-01-02"))
}

func decodeLogEvent(data string) []map[string]interface{} {
	decoded, err := b64.StdEncoding.DecodeString(data)
	FailOnError(err, "Cant decode log data")
	reader := bytes.NewReader([]byte(decoded))
	gzreader, err := gzip.NewReader(reader)
	FailOnError(err, "Cant unzip data")
	output, err := ioutil.ReadAll(gzreader)
	FailOnError(err, "Cant unzip data")
	result := string(output)

	str := strings.Split(result, "\n")
	fmt.Println("Split length:", len(str))
	var events []map[string]interface{}
	for _, sub := range str {
		var jsonData map[string]interface{}
		json.Unmarshal([]byte(sub), &jsonData)
		events = append(events, jsonData)
	}

	return events
}
