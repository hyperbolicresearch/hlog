build:
	go build -o bin/hlog main.go kafka.go logs.go

test:
	go test -race ./...