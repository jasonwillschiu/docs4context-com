# docs4context

**Local Context7 Client for AI Agents**

docs4context is a Model Context Protocol (MCP) server that provides AI agents with tools to download and search repository context documents from context7.com. It offers basic lookup and filtering capabilities for locally stored context documents, making repository documentation accessible to AI assistants without requiring full repository downloads.

**Currently only tested on macOS ARM (Apple Silicon) machines**
I'm personally using macOS ARM (darwin/arm64) and have only tested for the mcp working there so far.

## 🚀 Quick Start

### One-Line Installation
This installs into your current folder as a standalone executable.
```bash
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh
```

### Verify Installation
```bash
docs4context-com --version
```

### Configure Your AI Client
Add to your MCP client configuration (e.g., `opencode.json` for OpenCode):
```json
{
  "mcp": {
    "docs4context": {
      "type": "local",
      "command": ["docs4context-com"],
      "environment": {}
    }
  }
}
```

### First Use
Ask your AI assistant:
```
Download the context for the repository "mark3labs/mcp-go"
```

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Available Tools](#-available-tools)
- [Use Cases](#-use-cases)
- [How It Works](#-how-it-works)
- [Troubleshooting](#-troubleshooting)
- [Development](#-development)
- [Contributing](#-contributing)

## ✨ Features

### 🔍 **Context7 Document Access**
- Download pre-processed context documents from context7.com
- Search across locally stored context documents
- Cross-repository keyword analysis
- Topic-based content discovery with line number references

### 📊 **Local Context Management**
- Accurate token counting using GPT-4 compatible encoding
- Local storage of context7.com llms.txt files
- Metadata tracking for each repository
- Efficient local search and retrieval

### 🔧 **AI Agent Integration**
- Native MCP (Model Context Protocol) support
- Works with Claude Code, OpenCode, and other MCP clients
- Easy integration with existing AI workflows
- No manual configuration required

### 🌐 **Multi-Platform Support**
- macOS (Intel & Apple Silicon)
- Linux (x64 & ARM64)
- Windows (x64 & ARM64)
- Automatic updates and version management

## 🛠️ Installation

### Method 1: One-Line Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh
```

### Method 2: Specific Version
```bash
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh -s -- --version 0.1.2
```

### Method 3: Manual Installation
1. Download the appropriate binary from [GitHub Releases](https://github.com/jasonwillschiu/docs4context-com/releases)
2. Copy to `~/.local/bin/docs4context-com`
3. Make executable: `chmod +x ~/.local/bin/docs4context-com`
4. Add `~/.local/bin` to your PATH

### Method 4: Build from Source
```bash
git clone https://github.com/jasonwillschiu/docs4context-com
cd docs4context-com
go build -o docs4context-com .
```

## ⚙️ Configuration

### MCP Client Setup

#### For Claude Code
Use the CLI command to add the server:
```bash
# For local installation
claude mcp add docs4context -- docs4context-com

# Alternative: Remote server (if available)
# claude mcp add --transport sse docs4context https://your-server-url/sse
```

#### For Cursor
Add to your `mcp.json` configuration file:
```json
{
  "mcpServers": {
    "docs4context": {
      "command": "docs4context-com",
      "args": []
    }
  }
}
```

Alternative configurations for Cursor:
```json
// Using npx (if not globally installed)
{
  "mcpServers": {
    "docs4context": {
      "command": "npx",
      "args": ["-y", "docs4context-com"]
    }
  }
}

// Using bunx
{
  "mcpServers": {
    "docs4context": {
      "command": "bunx",
      "args": ["-y", "docs4context-com"]
    }
  }
}
```

#### For Windsurf
Add to your MCP configuration:
```json
{
  "mcpServers": {
    "docs4context": {
      "command": "docs4context-com",
      "args": []
    }
  }
}
```

Alternative for npx:
```json
{
  "mcpServers": {
    "docs4context": {
      "command": "npx",
      "args": ["-y", "docs4context-com"]
    }
  }
}
```

#### For OpenCode
Create or update `opencode.json`:
```json
{
  "mcp": {
    "docs4context": {
      "type": "local",
      "command": ["docs4context-com"],
      "environment": {}
    }
  }
}
```

#### For Other MCP Clients
The server uses stdio transport and works with any MCP-compatible client. General pattern:
```json
{
  "mcpServers": {
    "docs4context": {
      "command": "docs4context-com",
      "args": []
    }
  }
}
```

## 🛠️ Available Tools

### Document Management
- **`save_context_document`** - Downloads and saves repository context documents
  - Supports GitHub URLs or `username/repo` format
  - Includes accurate token counting
  - Adds metadata headers with creation time and source

### Search & Discovery
- **`search_titles`** - Find topics by title keywords
  - Optional repository filtering
  - Returns matching topics with line numbers

- **`search_content`** - Search descriptions and code content
  - Full-text search across repository content
  - Context-aware results with surrounding text

- **`get_topic_details`** - Extract complete topic information
  - Retrieve detailed content from specific line numbers
  - Includes surrounding context for better understanding

### Repository Management
- **`list_repositories`** - Show all available repositories
  - Displays metadata and topic counts
  - Helps discover what's already downloaded

- **`analyze_keywords`** - Cross-repository keyword analysis
  - Frequency analysis across all repositories
  - Relevance scoring for better search results

## 🎯 Use Cases

### Use Case 1: Learning a New Framework

**Scenario**: You want to understand how to implement authentication in a Go web framework.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AI Assistant  │    │  docs4context    │    │   context7.com  │
│                 │    │     Server       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         │ 1. "Download gin-gonic/gin"                     │
         ├───────────────────────►│                        │
         │                        │ 2. Fetch llms.txt      │
         │                        ├───────────────────────►│
         │                        │                        │
         │                        │ 3. Save locally        │
         │                        │                        │
         │ 4. "Search for auth"   │                        │
         ├───────────────────────►│                        │
         │                        │ 5. Search local content
         │                        │                        │
         │ 6. Results with line numbers                     │
         │◄───────────────────────┤                        │
         │                        │                        │
         │ 7. "Get details for auth middleware"             │
         ├───────────────────────►│                        │
         │                        │                        │
         │ 8. Complete topic details                       │
         │◄───────────────────────┤                        │
```

**Commands**:
1. `save_context_document` with GitHub URL: `gin-gonic/gin`
2. `search_titles` with query: `authentication`
3. `search_content` with query: `middleware auth`
4. `get_topic_details` for specific authentication examples

**Result**: Complete understanding of authentication patterns in Gin framework with specific code examples and implementation details.

### Use Case 2: Cross-Repository Code Pattern Analysis

**Scenario**: You're researching how different projects implement error handling patterns.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AI Assistant  │    │  docs4context    │    │   context7.com  │
│                 │    │     Server       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         │ 1. Download multiple repos                      │
         ├───────────────────────►│                        │
         │                        │ 2. Fetch llms.txt files│
         │                        ├───────────────────────►│
         │                        │                        │
         │ 3. "Analyze keyword: error"                     │
         ├───────────────────────►│                        │
         │                        │ 4. Search local files │
         │                        │                        │
         │ 5. Frequency & relevance scores                 │
         │◄───────────────────────┤                        │
         │                        │                        │
         │ 6. "Search error handling in Go projects"       │
         ├───────────────────────►│                        │
         │                        │                        │
         │ 7. Pattern examples from local context          │
         │◄───────────────────────┤                        │
```

**Commands**:
1. `save_context_document` for multiple repositories (e.g., `gin-gonic/gin`, `gorilla/mux`, `echo-labstack/echo`)
2. `analyze_keywords` with keyword: `error`
3. `search_content` with query: `error handling patterns`
4. `get_topic_details` for specific error handling implementations

**Result**: Comprehensive comparison of error handling approaches across different Go web frameworks with specific examples and best practices.

### Use Case 3: API Documentation Discovery

**Scenario**: You need to understand how to use a specific API or library feature.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AI Assistant  │    │  docs4context    │    │   Local Context │
│                 │    │     Server       │    │    Documents    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         │ 1. "What repos are available?"                  │
         ├───────────────────────►│                        │
         │                        │ 2. Check local storage│
         │                        ├───────────────────────►│
         │ 3. List all local repos                         │
         │◄───────────────────────┤                        │
         │                        │                        │
         │ 4. "Search for websocket implementation"        │
         ├───────────────────────►│                        │
         │                        │ 5. Search local files │
         │                        ├───────────────────────►│
         │ 6. WebSocket topics found                       │
         │◄───────────────────────┤                        │
         │                        │                        │
         │ 7. "Get complete WebSocket setup guide"         │
         ├───────────────────────►│                        │
         │                        │                        │
         │ 8. Detailed implementation guide                │
         │◄───────────────────────┤                        │
```

**Commands**:
1. `list_repositories` to see what's available
2. `search_titles` with query: `websocket`
3. `search_content` with query: `websocket implementation setup`
4. `get_topic_details` for complete WebSocket implementation examples

**Result**: Step-by-step guide for implementing WebSocket functionality with complete code examples and configuration details.

## 🔧 How It Works

### Context Document Flow
1. **Download**: Fetches pre-processed llms.txt files from context7.com
2. **Process**: Counts tokens using GPT-4 compatible encoding
3. **Store**: Saves locally with metadata headers including token count, date, and source
4. **Search**: Provides basic search and filtering across locally stored documents

### File Structure
```
llm-context/
├── username1/
│   ├── repo1/
│   │   └── llms.txt          # Context document with metadata
│   └── repo2/
│       └── llms.txt
└── username2/
    └── repo3/
        └── llms.txt
```

### Metadata Format
Each context document includes a header with:
```
# METADATA
# TOKEN_COUNT: 66551
# DATE_CREATED: 2025-06-26T10:30:45Z
# REPO: mark3labs/mcp-go
# SOURCE: https://context7.com/mark3labs/mcp-go/llms.txt
#
```

## 🐛 Troubleshooting

### Common Issues

#### Binary Not Found
```bash
# Check if binary is in PATH
which docs4context-com

# Add to PATH if needed
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Permission Denied
```bash
# Make binary executable
chmod +x ~/.local/bin/docs4context-com
```

#### MCP Client Not Connecting
1. Verify binary path in MCP configuration
2. Check that binary is executable
3. Review MCP client logs for connection errors
4. Test binary directly: `docs4context-com --version`

#### Context Download Fails
- Check internet connection
- Verify repository exists and is public
- Try with different repository format (`username/repo` vs full URL)

### Version Management
```bash
# Check current version
docs4context-com --version

# Check for updates
docs4context-com --check-updates

# Update to latest version
docs4context-com --update
```

### Debug Mode
Set environment variable for detailed logging:
```bash
DEBUG=true docs4context-com
```

## 🏗️ Development

### Prerequisites
- Go 1.24.2 or later
- Bun (for build scripts)

### Build Commands
```bash
# Development build
go build -o docs4context-com .

# Cross-platform builds
bun run cicd.js --build

# Build with version info
bun run cicd.js --build --commit --tag --release
```

### Testing
```bash
# Install dependencies
go mod tidy

# Test local build
./docs4context-com --version
```

### Architecture
- **MCP Server**: Built with `github.com/mark3labs/mcp-go`
- **Token Counting**: Uses `github.com/pkoukk/tiktoken-go`
- **Transport**: stdio-based communication
- **Storage**: Local file system with organized directory structure

## 🤝 Contributing

### Reporting Issues
Report bugs and feature requests at: https://github.com/jasonwillschiu/docs4context-com/issues

### Development Process
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

### Code Style
- Follow Go conventions
- Add tests for new features
- Update documentation
- Ensure cross-platform compatibility

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built on [mcp-go](https://github.com/mark3labs/mcp-go) framework
- Powered by [context7.com](https://context7.com) for repository context
- Token counting via [tiktoken-go](https://github.com/pkoukk/tiktoken-go)

---

**Made with ❤️ for AI-assisted development**
