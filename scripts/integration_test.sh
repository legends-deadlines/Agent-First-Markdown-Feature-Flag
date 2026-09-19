#!/bin/bash
# Интеграционный тест MDFLAG
# Тестирует полный цикл: создание → верификация → роллаут → верификация

set -e

echo "=== MDFLAG Integration Test ==="
echo ""

# Настройка
TEST_DIR=$(mktemp -d)
FLAGS_DIR="$TEST_DIR/.mdflag"
BINARY="./mdflag"

echo "Test directory: $TEST_DIR"
echo "Flags directory: $FLAGS_DIR"
echo ""

# Очистка при выходе
cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Проверяем, что бинарник существует
if [ ! -f "$BINARY" ]; then
    echo "ERROR: Binary $BINARY not found. Run 'go build -o mdflag ./cmd/mdflag/' first."
    exit 1
fi

# ==========================================
# Тест 1: Создание флага
# ==========================================
echo "--- Test 1: Create flag ---"

$BINARY create \
    --name test-flag \
    --percentage 0 \
    --hypothesis "Test hypothesis" \
    --metrics "metric1,metric2" \
    --author "test-agent" \
    --dir "$FLAGS_DIR"

if [ ! -f "$FLAGS_DIR/test-flag.md" ]; then
    echo "FAIL: Flag file not created"
    exit 1
fi
echo "PASS: Flag file created"
echo ""

# ==========================================
# Тест 2: Проверка содержимого файла
# ==========================================
echo "--- Test 2: Check file content ---"

if ! grep -q "name: test-flag" "$FLAGS_DIR/test-flag.md"; then
    echo "FAIL: Flag name not found in file"
    exit 1
fi

if ! grep -q "percentage: 0" "$FLAGS_DIR/test-flag.md"; then
    echo "FAIL: Percentage not found in file"
    exit 1
fi

if ! grep -q "agent_section_hash:" "$FLAGS_DIR/test-flag.md"; then
    echo "FAIL: Agent section hash not found"
    exit 1
fi

if ! grep -q "human_section_hash:" "$FLAGS_DIR/test-flag.md"; then
    echo "FAIL: Human section hash not found"
    exit 1
fi

echo "PASS: File content is correct"
echo ""

# ==========================================
# Тест 3: Верификация целостности
# ==========================================
echo "--- Test 3: Verify integrity ---"

$BINARY verify --dir "$FLAGS_DIR"

echo "PASS: Verification passed"
echo ""

# ==========================================
# Тест 4: Роллаут (изменение процента)
# ==========================================
echo "--- Test 4: Rollout to 25% ---"

$BINARY rollout --name test-flag --percentage 25 --dir "$FLAGS_DIR"

if ! grep -q "percentage: 25" "$FLAGS_DIR/test-flag.md"; then
    echo "FAIL: Percentage not updated"
    exit 1
fi

echo "PASS: Percentage updated to 25%"
echo ""

# ==========================================
# Тест 5: Верификация после роллаута
# ==========================================
echo "--- Test 5: Verify after rollout ---"

$BINARY verify --dir "$FLAGS_DIR"

echo "PASS: Verification passed after rollout"
echo ""

# ==========================================
# Тест 6: Список флагов
# ==========================================
echo "--- Test 6: List flags ---"

$BINARY list --dir "$FLAGS_DIR"

echo "PASS: List command works"
echo ""

# ==========================================
# Тест 7: Дубликат флага должен отклоняться
# ==========================================
echo "--- Test 7: Duplicate flag rejection ---"

if $BINARY create \
    --name test-flag \
    --percentage 0 \
    --hypothesis "Duplicate" \
    --dir "$FLAGS_DIR" 2>/dev/null; then
    echo "FAIL: Duplicate flag should be rejected"
    exit 1
fi

echo "PASS: Duplicate flag rejected"
echo ""

# ==========================================
# Тест 8: Невалидный процент должен отклоняться
# ==========================================
echo "--- Test 8: Invalid percentage rejection ---"

if $BINARY create \
    --name invalid-flag \
    --percentage 150 \
    --hypothesis "Invalid" \
    --dir "$FLAGS_DIR" 2>/dev/null; then
    echo "FAIL: Invalid percentage should be rejected"
    exit 1
fi

echo "PASS: Invalid percentage rejected"
echo ""

# ==========================================
# Тест 9: Обнаружение подмены файла
# ==========================================
echo "--- Test 9: Tamper detection ---"

# Создаём второй флаг
$BINARY create \
    --name tamper-test \
    --percentage 5 \
    --hypothesis "Tamper test" \
    --dir "$FLAGS_DIR"

# Подменяем процент напрямую в файле (без инструмента)
sed -i 's/percentage: 5/percentage: 100/' "$FLAGS_DIR/tamper-test.md"

# Верификация должна обнаружить подмену
if $BINARY verify --dir "$FLAGS_DIR" 2>/dev/null; then
    echo "FAIL: Tampered file should be detected"
    exit 1
fi

echo "PASS: Tampered file detected"
echo ""

# ==========================================
# Итог
# ==========================================
echo "=== All Integration Tests Passed ==="
