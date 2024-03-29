build:
	go build -o bin/hlog ./cmd/hlog
	go build -o bin/randomproducer ./cmd/randomproducer

test:
	go test -race ./...