# Calculator MCP Server - opencode Integration

## ✅ Setup Complete!

Your calculator MCP server has been successfully created and configured for opencode integration.

## Files Created

### 1. Calculator MCP Server
- **Location**: `calculator-mcp/`
- **Binary**: `calculator-mcp/calculator-mcp`
- **Source**: `calculator-mcp/main.go`
- **Documentation**: `calculator-mcp/README.md`

### 2. opencode Configuration
- **File**: `opencode.json`
- **Purpose**: Configures opencode to use the calculator MCP server

### 3. Integration Test
- **File**: `test-calculator-integration.sh`
- **Purpose**: Verifies the setup is working correctly

## Available Calculator Tools

The MCP server provides these tools to opencode:

1. **add** - Add two numbers together
2. **subtract** - Subtract the second number from the first
3. **multiply** - Multiply two numbers together
4. **divide** - Divide the first number by the second (with zero-division protection)
5. **power** - Raise the first number to the power of the second
6. **sqrt** - Calculate the square root of a number (with negative number protection)

## How to Use

### 1. Start opencode
```bash
opencode
```

### 2. Ask for calculations
Once opencode is running, you can ask it to perform calculations using natural language:

- "Can you calculate 15 + 27?"
- "What's 144 divided by 12?"
- "Calculate the square root of 64"
- "What's 2 to the power of 10?"
- "Multiply 7 by 8"

### 3. Monitor logs
The calculator server logs all operations to stderr, so you can see when tools are being used:

```
[Calculator MCP] 2025/06/25 12:01:40 main.go:20: Starting Calculator MCP Server v1.0.0
[Calculator MCP] 2025/06/25 12:01:45 main.go:65: ADD tool called
[Calculator MCP] 2025/06/25 12:01:45 main.go:77: ADD tool: 15 + 27 = 42
```

## Configuration Details

The `opencode.json` configuration tells opencode:
- **Server name**: "calculator"
- **Type**: "local" (runs on your machine)
- **Command**: Points to the calculator binary
- **Environment**: No special environment variables needed

## Troubleshooting

### If tools aren't available:
1. Restart opencode to reload the configuration
2. Check that `calculator-mcp/calculator-mcp` exists and is executable
3. Verify the path in `opencode.json` is correct

### If calculations fail:
1. Check the calculator server logs for error messages
2. Ensure you're providing valid numbers
3. Watch for division by zero or negative square root errors

## Next Steps

You can now:
1. **Test the integration** by asking opencode to perform calculations
2. **Extend the calculator** by adding more mathematical functions
3. **Create additional MCP servers** for other tools and services
4. **Share your setup** with others who want to use MCP servers with opencode

## Success! 🎉

Your calculator MCP server is ready to use with opencode. The integration provides a seamless way for the AI to perform mathematical calculations using your custom tools.