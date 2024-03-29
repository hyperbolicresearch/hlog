build:
	go build -o bin/hlog ./cmd/hlog
	go build -o bin/hlog_producer ./cmd/hlog_producer

test:
	go test -race ./...