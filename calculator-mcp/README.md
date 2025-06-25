# Calculator MCP Server

A basic calculator MCP (Model Context Protocol) server built with Go using the [mcp-go](https://github.com/mark3labs/mcp-go) framework.

## Features

This MCP server provides the following calculator tools:

- **add**: Add two numbers together
- **subtract**: Subtract the second number from the first number  
- **multiply**: Multiply two numbers together
- **divide**: Divide the first number by the second number (with zero-division protection)
- **power**: Raise the first number to the power of the second number
- **sqrt**: Calculate the square root of a number (with negative number protection)

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the server:
   ```bash
   go build -o calculator-mcp main.go
   ```

## Usage

### Running the Server

The server uses stdio transport by default:

```bash
./calculator-mcp
```

### Logging

The server includes comprehensive logging to stderr for debugging and monitoring:

- **Startup logs**: Server initialization and tool registration
- **Tool usage logs**: Each time a calculator tool is called with input/output values
- **Error logs**: Detailed error information for troubleshooting

Example log output:
```
[Calculator MCP] 2025/06/25 12:01:40 main.go:20: Starting Calculator MCP Server v1.0.0
[Calculator MCP] 2025/06/25 12:01:40 main.go:31: Registering calculator tools...
[Calculator MCP] 2025/06/25 12:01:40 main.go:38: All tools registered successfully
[Calculator MCP] 2025/06/25 12:01:40 main.go:41: Starting stdio server...
[Calculator MCP] 2025/06/25 12:01:45 main.go:65: ADD tool called
[Calculator MCP] 2025/06/25 12:01:45 main.go:77: ADD tool: 5 + 3 = 8
```

The logs are sent to stderr to avoid interfering with the MCP protocol communication on stdout.

### Tool Examples

#### Addition
```json
{
  "tool": "add",
  "arguments": {
    "a": 5,
    "b": 3
  }
}
```
Result: `8`

#### Subtraction
```json
{
  "tool": "subtract", 
  "arguments": {
    "a": 10,
    "b": 4
  }
}
```
Result: `6`

#### Multiplication
```json
{
  "tool": "multiply",
  "arguments": {
    "a": 6,
    "b": 7
  }
}
```
Result: `42`

#### Division
```json
{
  "tool": "divide",
  "arguments": {
    "a": 15,
    "b": 3
  }
}
```
Result: `5`

#### Power/Exponentiation
```json
{
  "tool": "power",
  "arguments": {
    "base": 2,
    "exponent": 8
  }
}
```
Result: `256`

#### Square Root
```json
{
  "tool": "sqrt",
  "arguments": {
    "number": 16
  }
}
```
Result: `4`

## Error Handling

The server includes proper error handling for:

- **Division by zero**: Returns error message "cannot divide by zero"
- **Square root of negative numbers**: Returns error message "cannot calculate square root of negative number"
- **Invalid mathematical results**: Handles NaN and infinite results in power operations
- **Missing or invalid arguments**: Validates all required parameters

## Integration with MCP Clients

This server can be integrated with any MCP-compatible client, such as:

- **opencode**: Add to your `opencode.json` configuration
- **Claude Desktop**: Add to your MCP configuration
- Other LLM applications that support MCP
- Custom MCP clients

### Integration with opencode

To use this calculator with opencode, create or update your `opencode.json` file:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "calculator": {
      "type": "local",
      "command": ["./calculator-mcp/calculator-mcp"],
      "environment": {}
    }
  }
}
```

**Note**: Make sure the path in the `command` array points to the correct location of your `calculator-mcp` binary.

After adding this configuration:
1. Restart opencode if it's already running
2. The calculator tools will be automatically available to the LLM
3. You can ask opencode to perform calculations using natural language

Example usage in opencode:
- "Can you calculate 15 divided by 3?"
- "What's the square root of 144?"
- "Calculate 2 to the power of 8"

### Integration with Claude Desktop

For Claude Desktop, add this to your MCP configuration file:

```json
{
  "mcpServers": {
    "calculator": {
      "command": "/path/to/calculator-mcp",
      "args": []
    }
  }
}
```

## Development

### Project Structure

```
calculator-mcp/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # This documentation
```

### Adding New Tools

To add new calculator functions:

1. Create a new tool definition using `mcp.NewTool()`
2. Define the required parameters with appropriate types and descriptions
3. Implement the handler function with proper error handling
4. Add the tool to the server using `s.AddTool()`

### Dependencies

- [mcp-go](https://github.com/mark3labs/mcp-go): Go implementation of the Model Context Protocol

## License

This project is provided as an example implementation. Feel free to modify and use as needed.