#!/bin/bash

echo "=== Calculator MCP Server Integration Test ==="
echo

# Check if opencode.json exists
if [ -f "opencode.json" ]; then
    echo "✓ opencode.json configuration file found"
    echo "Configuration:"
    cat opencode.json | jq '.mcp.calculator' 2>/dev/null || echo "  (JSON parsing not available, showing raw content)"
    echo
else
    echo "✗ opencode.json not found"
    exit 1
fi

# Check if calculator binary exists
if [ -f "calculator-mcp/calculator-mcp" ]; then
    echo "✓ Calculator MCP binary found"
    ls -la calculator-mcp/calculator-mcp
    echo
else
    echo "✗ Calculator MCP binary not found"
    exit 1
fi

# Test calculator server startup
echo "Testing calculator server startup..."
cd calculator-mcp
./calculator-mcp &
CALC_PID=$!
sleep 2

if kill -0 $CALC_PID 2>/dev/null; then
    echo "✓ Calculator server started successfully (PID: $CALC_PID)"
    kill $CALC_PID 2>/dev/null
    wait $CALC_PID 2>/dev/null
    echo "✓ Calculator server stopped cleanly"
else
    echo "✗ Calculator server failed to start or exited early"
fi

cd ..
echo
echo "=== Integration Test Complete ==="
echo
echo "To use with opencode:"
echo "1. Make sure opencode.json is in your project directory"
echo "2. Restart opencode if it's already running"
echo "3. Ask opencode to perform calculations like:"
echo "   - 'Calculate 5 + 3'"
echo "   - 'What's the square root of 16?'"
echo "   - 'Divide 20 by 4'"