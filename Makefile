build:
	go build -o bin/hlog main.go kafka.go logs.go
	go build -o bin/hlog_producer ./cmd/hlog_producer

test:
	go test -race ./...