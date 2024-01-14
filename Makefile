build:
	go build -o bin/hlog cmd/hlog/main.go

test:
	go test -race ./...