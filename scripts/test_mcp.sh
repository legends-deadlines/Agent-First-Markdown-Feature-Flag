#!/bin/bash
# Ручное тестирование MCP-сервера

set -e

BINARY="./mdflag"
FLAGS_DIR=$(mktemp -d)

echo "=== MCP Server Manual Test ==="
echo "Flags directory: $FLAGS_DIR"
echo ""

cleanup() {
    rm -rf "$FLAGS_DIR"
}
trap cleanup EXIT

# Проверяем бинарник
if [ ! -f "$BINARY" ]; then
    echo "ERROR: Binary not found. Build first."
    exit 1
fi

echo "Testing MCP server initialization..."
echo ""

# Отправляем запрос инициализации
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"0.1.0"}}}' | timeout 5 $BINARY serve --dir "$FLAGS_DIR" || true

echo ""
echo "If you see a JSON response above, MCP server initialization works."
echo ""

echo "Testing tools/list..."
echo ""

# Запрашиваем список инструментов
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | timeout 5 $BINARY serve --dir "$FLAGS_DIR" || true

echo ""
echo "If you see mdflag_create, mdflag_list, mdflag_verify in the response, tools are registered."
echo ""

echo "Testing mdflag_create tool..."
echo ""

# Создаём флаг через MCP
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"mdflag_create","arguments":{"name":"mcp-test-flag","hypothesis":"MCP integration test","percentage":0,"author":"mcp-test"}}}' | timeout 5 $BINARY serve --dir "$FLAGS_DIR" || true

echo ""

# Проверяем, что файл создан
if [ -f "$FLAGS_DIR/mcp-test-flag.md" ]; then
    echo "PASS: Flag file created via MCP"
else
    echo "FAIL: Flag file not created"
fi

echo ""
echo "Testing mdflag_list tool..."
echo ""

echo '{"jsonrpc":"2.0","id":4,"method":"tools/list","params":{}}' | timeout 5 $BINARY serve --dir "$FLAGS_DIR" || true

echo ""
echo "=== MCP Tests Complete ==="
