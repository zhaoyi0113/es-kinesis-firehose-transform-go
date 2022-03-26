build:
	GOARCH=darwin GOARCH=amd64 go build -o dist/api

buildimage:
	docker build --platform=linux/amd64 -t zhaoyi0113/es-kinesis-firehose-transform-go .

publishimage:
	docker push zhaoyi0113/es-kinesis-firehose-transform-go

unittest:
	go test -v ./...
