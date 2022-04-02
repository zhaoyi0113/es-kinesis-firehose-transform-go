package internal

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed testData/logTestData.json
var logTestData string

func TestLogProcess(t *testing.T) {
	var logEvents []LogEvent
	err := json.Unmarshal([]byte(logTestData), &logEvents)
	FailOnError(err, "Cant parse log test data")

	for _, logEvent := range logEvents {
		record := LogEventRecord{
			Req: struct {
				Body LogEvent `json:"body"`
			}{Body: logEvent},
		}
		ProcessLogs(record, "logs")
	}
}

func TestGetIndexName(t *testing.T) {
	name := getIndexName("logs")
	now := time.Now()
	expected := fmt.Sprintf("%s-%s", "aws-logs", now.Format("2006-01-02"))
	assert.Equal(t, name, expected)
}

func TestDecodeEvent(t *testing.T) {

}
