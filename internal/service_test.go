package internal

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed testData/logTestData.json
var logTestData string

var es = CreateESClient()

func TestLogProcess(t *testing.T) {
	var logEvents []LogEvent
	err := json.Unmarshal([]byte(logTestData), &logEvents)
	FailOnError(err, "Cant parse log test data")

	var record LogEventRecord
	record.Body = logEvents[0]
	response := ProcessLogs(record, "logs")
	assert.Equal(t, response.RequestId, logEvents[0].RequestId)

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"@message": "7c883b5cc430: Download complete",
			},
		},
	}
	json.NewEncoder(&buf).Encode(query)
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(getIndexName("logs")),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	var r map[string]interface{}
	json.NewDecoder(res.Body).Decode(&r)
	fmt.Println("Get es response:", r)
	hits := r["hits"]
	total := hits.(map[string]interface{})["total"]
	valueStr := total.(map[string]interface{})["value"]
	assert.Equal(t, reflect.ValueOf(valueStr).Float() > 0, true)
}

func TestGetIndexName(t *testing.T) {
	name := getIndexName("logs")
	now := time.Now()
	expected := fmt.Sprintf("%s-%s", "aws-logs", now.Format("2006-01-02"))
	assert.Equal(t, name, expected)
}
