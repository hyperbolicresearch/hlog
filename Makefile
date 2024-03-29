build:
	go build -o bin/hlog ./cmd/hlog
	go build -o bin/hlog_producer ./cmd/hlog_producer
	go build -o bin/hlog_livetail ./cmd/hlog_livetail

test:
	go test -race ./...