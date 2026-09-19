#!/bin/bash
# Нагрузочное тестирование: много флагов, параллельные запросы

set -e

BINARY="./mdflag"
FLAGS_DIR=$(mktemp -d)
NUM_FLAGS=100

echo "=== Benchmark Test ==="
echo "Creating $NUM_FLAGS flags..."

cleanup() {
    rm -rf "$FLAGS_DIR"
}
trap cleanup EXIT

# Создаём много флагов
for i in $(seq 1 $NUM_FLAGS); do
    $BINARY create \
        --name "benchmark-flag-$i" \
        --percentage 50 \
        --hypothesis "Benchmark test $i" \
        --dir "$FLAGS_DIR" > /dev/null
done

echo "Created $NUM_FLAGS flags."
echo ""

# Тестируем время верификации
echo "Verifying $NUM_FLAGS flags..."
time $BINARY verify --dir "$FLAGS_DIR" > /dev/null

echo ""

# Тестируем время списка
echo "Listing $NUM_FLAGS flags..."
time $BINARY list --dir "$FLAGS_DIR" > /dev/null

echo ""
echo "=== Benchmark Complete ==="
