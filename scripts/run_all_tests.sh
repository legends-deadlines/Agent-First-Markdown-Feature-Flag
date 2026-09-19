#!/bin/bash
set -e

BINARY="./mdflag"

if ! command -v go &> /dev/null; then
    echo "ОШИБКА: команда 'go' не найдена."
    echo "Убедитесь, что Go установлен и добавлен в PATH."
    echo "Пример: export PATH=\$PATH:/usr/local/go/bin"
    exit 1
fi

echo "=== Запуск всех тестов ==="

echo "1. Сборка бинарника..."
go build -o "$BINARY" ./cmd/mdflag/

echo "2. Юнит-тесты и покрытие..."
go test ./... -v -cover

echo "3. Интеграционные тесты..."
chmod +x scripts/integration_test.sh
./scripts/integration_test.sh

echo "4. Тесты MCP-сервера..."
chmod +x scripts/test_mcp.sh
./scripts/test_mcp.sh

echo "5. Нагрузочные тесты..."
chmod +x scripts/benchmark_test.sh
./scripts/benchmark_test.sh

echo "=== Все тесты завершены ==="
