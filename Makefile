build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

run:
	go run cmd/hexlet-path-size/main.go

# clean:
# 	rm -rf bin/*

deps:
	go mod download

lint:
	golangci-lint run ./...

test:
	go test -v ./...
# 	go test ./... -v | grep -v "no test files"