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
	Body LogEvent `json:"body"`
}

type LogEvent struct {
	Records []struct {
		Data string `json:"data"`
	} `json:"records"`
	Timestamp int    `json:"timestamp"`
	RequestId string `json:"requestId"`
}

type LogDataEvent struct {
	Id        string      `json:"id"`
	Timestamp json.Number `json:"timestamp"`
	Message   string      `json:"message"`
}

type LogData struct {
	MessageType         string         `json:"messageType"`
	Owner               string         `json:"owner"`
	LogGroup            string         `json:"logGroup"`
	LogStream           string         `json:"logStream"`
	SubscriptionFilters []string       `json:"subscriptionFilters"`
	LogEvents           []LogDataEvent `json:"logEvents"`
}

type ESDoc struct {
	Message string `json:"@message"`
}

type Response struct {
	Version     string `json:"X-Amz-Firehose-Protocol-Version"`
	RequestId   string `json:"X-Amz-Firehose-Request-Id"`
	ContentType string `json:"Content-Type"`
}

func ProcessLogs(event LogEventRecord, eventType string) Response {
	indexName := getIndexName("logs")
	CreateIndex(indexName)
	var bulkIndex BulkIndex
	bulkIndex.Index.Index = indexName
	bulkIndexStr, _ := json.Marshal(bulkIndex)
	bulkDocs := ""
	for _, record := range event.Body.Records {
		if eventType == "logs" {
			logs := decodeLogEvent(record.Data)
			for _, log := range logs {
				for _, logEvent := range log.LogEvents {
					fmt.Println("process log event:", logEvent.Id, logEvent.Timestamp)
					if len(logEvent.Message) > 0 {
						var jsonData map[string]string
						err := json.Unmarshal([]byte(logEvent.Message), &jsonData)
						esDoc := make(map[string]string)
						if err != nil {
							esDoc["@message"] = logEvent.Message
						} else if v, found := jsonData["@message"]; found {
							esDoc["@message"] = v
						} else {
							esDoc["@message"] = logEvent.Message
						}
						esDoc["messageType"] = log.MessageType
						esDoc["owner"] = log.Owner
						esDoc["logGroup"] = log.LogGroup
						esDoc["logStream"] = log.LogStream
						esDoc["id"] = logEvent.Id
						esDoc["timestamp"] = string(logEvent.Timestamp)
						esJsonDoc, _ := json.Marshal(esDoc)
						bulkDocs += string(bulkIndexStr) + "\n"
						bulkDocs += string(esJsonDoc) + "\n"
					}
				}
			}
		}
	}
	BulkInsert(indexName, bulkDocs)
	var response Response
	response.ContentType = "application/json"
	response.Version = "1.0"
	response.RequestId = event.Body.RequestId
	return response
}

func getIndexName(indexType string) string {
	now := time.Now()
	return fmt.Sprintf("aws-%s-%s", indexType, now.Format("2006-01-02"))
}

func decodeLogEvent(data string) []LogData {
	decoded, err := b64.StdEncoding.DecodeString(data)
	FailOnError(err, "Cant decode log data")
	reader := bytes.NewReader([]byte(decoded))
	gzreader, err := gzip.NewReader(reader)
	FailOnError(err, "Cant unzip data")
	output, err := ioutil.ReadAll(gzreader)
	FailOnError(err, "Cant unzip data")
	result := string(output)

	str := strings.Split(result, "\n")
	fmt.Println("Original:", str)
	var events []LogData
	for _, sub := range str {
		var logData LogData
		json.Unmarshal([]byte(sub), &logData)
		events = append(events, logData)
	}
	return events
}
