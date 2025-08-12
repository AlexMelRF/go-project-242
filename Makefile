# Сборка бинарного файла
build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

# Запуск программы
run:
	go run cmd/hexlet-path-size/main.go

# Очистка собранных файлов
clean:
	rm -rf bin/*

# Установка зависимостей
deps:
	go mod download

# Запуск линтера
lint:
	golangci-lint run ./...