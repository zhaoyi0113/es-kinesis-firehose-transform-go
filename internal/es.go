package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type IndexMapping struct {
	Properties struct {
		Timestamp struct {
			Type string `json:"type"`
		} `json:"timestamp"`
	} `json:"properties"`
}

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
		FailOnError(err, "Failed to check index exists")
		return
	}
	res, err = client.Indices.Create(indexName)
	FailOnError(err, "Failed to create index "+indexName)
	failOnEsError(res)
	fmt.Println("Create index success:" + indexName)
	var mapping IndexMapping
	mapping.Properties.Timestamp.Type = "date"

	mappingJson, _ := json.Marshal(mapping)
	res, err = client.Indices.PutMapping([]string{indexName}, strings.NewReader(string(mappingJson)))
	FailOnError(err, "Failed to update index mapping")
	failOnEsError(res)
}

func failOnEsError(res *esapi.Response) {
	if res.IsError() {
		log.Println(res.String())
		panic(errors.New(res.String()))
	}
}
