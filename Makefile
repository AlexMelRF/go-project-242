build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

run:
	go run cmd/hexlet-path-size/main.go

clean:
	rm -rf bin/*

# lint:
# 	golangci-lint run ./..

test:
	go test -v ./...

install: build
	sudo cp bin/hexlet-path-size /usr/local/bin/hexlet-path-size

uninstall: 
	sudo rm /usr/local/bin/hexlet-path-size