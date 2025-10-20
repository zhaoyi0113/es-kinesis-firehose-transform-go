package internal

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func ProcessLogs(event LogEvent, eventType string) Response {

	indexName := getIndexName(eventType)
	CreateIndex(indexName)
	var bulkIndex BulkIndex
	bulkIndex.Index.Index = indexName
	bulkIndexStr, _ := json.Marshal(bulkIndex)
	var bulkDocs bytes.Buffer
	for _, record := range event.Records {
		if eventType == "logs" {
			logs := decodeLogEvent(record.Data)
			for _, log := range logs {
				for _, logEvent := range log.LogEvents {
					if len(logEvent.Message) > 0 {
						var jsonData map[string]string
						err := json.Unmarshal([]byte(logEvent.Message), &jsonData)
						esDoc := make(map[string]interface{})
						if err != nil {
							esDoc["@message"] = logEvent.Message
						} else if v, found := jsonData["@message"]; found {
							esDoc["@message"] = v
							esDoc["@componentName"] = jsonData["@componentName"]
							esDoc["@partName"] = jsonData["@partName"]
							esDoc["@region"] = jsonData["@region"]
							esDoc["@lambdaName"] = jsonData["@lambdaName"]
							esDoc["@level"] = jsonData["@level"]
							esDoc["@aggregateId"] = jsonData["@aggregateId"]
							esDoc["@requestId"] = jsonData["@requestId"]
							esDoc["@stage"] = jsonData["@stage"]
							esDoc["level"] = jsonData["level"]
						} else {
							esDoc["@message"] = logEvent.Message
						}
						esDoc["messageType"] = log.MessageType
						esDoc["owner"] = log.Owner
						esDoc["logGroup"] = log.LogGroup
						esDoc["logStream"] = log.LogStream
						esDoc["id"] = logEvent.Id
						esDoc["timestamp"] = string(logEvent.Timestamp)
						esDoc["timestampText"] = string(logEvent.Timestamp)
						var messageJsonData map[string]interface{}
						str := fmt.Sprintf("%v", esDoc["@message"])
						err = json.Unmarshal([]byte(str), &messageJsonData)

						if err == nil {
							v, exist := messageJsonData["@message"]
							if exist {
								esDoc["@message"] = v
								esDoc["@componentName"] = messageJsonData["@componentName"]
								esDoc["@partName"] = messageJsonData["@partName"]
								esDoc["@region"] = messageJsonData["@region"]
								esDoc["@lambdaName"] = messageJsonData["@lambdaName"]
								esDoc["@level"] = messageJsonData["@level"]
								esDoc["@aggregateId"] = messageJsonData["@aggregateId"]
								esDoc["@requestId"] = messageJsonData["@requestId"]
								esDoc["@stage"] = messageJsonData["@stage"]
								esDoc["level"] = messageJsonData["level"]
							}
						}

						esJsonDoc, _ := json.Marshal(esDoc)
						bulkDocs.Write(bulkIndexStr)
						bulkDocs.WriteString("\n")
						bulkDocs.Write(esJsonDoc)
						bulkDocs.WriteString("\n")
					}
				}
			}
		}
	}
	if bulkDocs.Len() > 0 {
		BulkInsert(indexName, bulkDocs.String())
	} else {
		log.Println("Not docs found for this event")
	}
	var response Response
	response.ContentType = "application/json"
	response.Version = "1.0"
	response.RequestId = event.RequestId
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
	output, err := io.ReadAll(gzreader)
	FailOnError(err, "Cant unzip data")
	result := string(output)
	gzreader.Close()

	str := strings.Split(result, "\n")
	var events []LogData
	for _, sub := range str {
		var logData LogData
		json.Unmarshal([]byte(sub), &logData)
		events = append(events, logData)
	}
	return events
}
