package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var client *elasticsearch.Client = nil

func getEsHost() string {
	esHost := os.Getenv("ES_HOST")
	if len(esHost) == 0 {
		return "http://localhost:9200"
	}
	return esHost
}

func createESClient() *elasticsearch.Client {
	esHost := getEsHost()
	cfg := elasticsearch.Config{
		Addresses: []string{
			esHost,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	FailOnError(err, "Failed to create elasticsearch connection.")
	return es
}

func CreateIndex(indexName string) {
	if client == nil {
		client = createESClient()
	}
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil || res.StatusCode == 200 {
		return
	}
	res, err = client.Indices.Create(indexName)
	FailOnError(err, "Failed to create index "+indexName)
	if res.IsError() {
		FailOnError(errors.New(res.String()), res.String())
	}
	fmt.Println("Create index success:" + indexName)
}
