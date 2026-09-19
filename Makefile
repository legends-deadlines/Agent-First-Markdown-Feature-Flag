.PHONY: all build test test-unit test-integration test-mcp test-benchmark clean install lint fmt coverage cross-build

BINARY_NAME=mdflag
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

all: lint test build

# Сборка бинарника для текущей платформы
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/mdflag/

# Запуск всех тестов
test: test-unit test-integration test-mcp test-benchmark

# Юнит-тесты с покрытием
test-unit:
	go test ./... -v -cover

# Интеграционные тесты
test-integration: build
	chmod +x scripts/integration_test.sh
	./scripts/integration_test.sh

# Тесты MCP-сервера
test-mcp: build
	chmod +x scripts/test_mcp.sh
	./scripts/test_mcp.sh

# Нагрузочные тесты
test-benchmark: build
	chmod +x scripts/benchmark_test.sh
	./scripts/benchmark_test.sh

# Проверка линтером
lint:
	go vet ./...

# Проверка форматирования
fmt:
	gofmt -w .
	gofmt -l .

# Отчёт о покрытии
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total

# Кроссплатформенная сборка
cross-build:
	@echo "Building for Linux amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/mdflag-linux-amd64 ./cmd/mdflag/
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/mdflag-linux-arm64 ./cmd/mdflag/
	@echo "Building for macOS amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/mdflag-darwin-amd64 ./cmd/mdflag/
	@echo "Building for macOS arm64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/mdflag-darwin-arm64 ./cmd/mdflag/
	@echo "Building for Windows amd64..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/mdflag-windows-amd64.exe ./cmd/mdflag/
	@echo "Generating checksums..."
	cd dist && sha256sum * > checksums.txt
	@echo "Done. Binaries in ./dist/"

# Установка в систему
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# Очистка
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html
