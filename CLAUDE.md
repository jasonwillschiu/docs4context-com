# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a collection of MCP (Model Context Protocol) servers providing tools for LLM agents:

1. **Main docs4context server** (`main.go`) - A local client for context7.com that downloads and provides basic search/filtering tools for repository context documents (llms.txt files)
2. **Calculator MCP server** (`calculator-mcp/main.go`) - Provides arithmetic operations

Both servers use the `github.com/mark3labs/mcp-go` framework with stdio transport.

## Development Commands

### Build
```bash
# Build for current platform (development version)
go build -o docs4context-com .

# Build for all platforms using cicd script
bun run cicd.js --build

# Build with version information (production)
bun run cicd.js --build --commit --tag --release

# Build and run locally
bun run cicd.js --mode build
bun run cicd.js --mode dev
```

### Cross-Platform Distribution
```bash
# Build binaries for all platforms
bun run cicd.js --build

# Creates binaries in bin/ directory:
# - docs4context-com-darwin-amd64
# - docs4context-com-darwin-arm64  
# - docs4context-com-linux-amd64
# - docs4context-com-linux-arm64
# - docs4context-com-windows-amd64.exe
# - docs4context-com-windows-arm64.exe
```

### Installation & Upgrades
```bash
# One-line install (latest release)
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh

# Install specific version
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh -s -- --version 0.1.2

# Check current version
docs4context-com --version

# Check for updates
docs4context-com --check-updates

# Auto-update to latest version
docs4context-com --update

# Manual install
# 1. Download appropriate binary from GitHub releases
# 2. Copy to ~/.local/bin/docs4context-com
# 3. Make executable: chmod +x ~/.local/bin/docs4context-com
```

### Test
```bash
# Install/update dependencies
go mod tidy

# Test local build
./docs4context-com
```

### Run
```bash
# Run main docs4context server
./docs4context-com

# Or if installed globally
docs4context-com
```

## Architecture

### MCP Server Pattern
- Both servers follow identical patterns using `mcp-go`
- stdio-based communication (logs go to stderr, never stdout)
- Tool registration with parameter validation
- Comprehensive error handling with fallbacks

### File Structure
- `main.go` - Main docs4context MCP server with search tools
- `calculator-mcp/` - Separate calculator MCP server with own module
- `llm-context/` - Downloaded context documents stored as `[username]/[repo]/llms.txt` with metadata headers
- `opencode.json` - MCP client configuration for opencode integration

### Key Dependencies
- Primary: `github.com/mark3labs/mcp-go v0.32.0`
- Token counting: `github.com/pkoukk/tiktoken-go v0.1.7`
- Go version: 1.24.2

## Distribution Strategy

This MCP server is distributed as **pre-built binaries** for easy installation:

### Platform Support
- **macOS**: Intel (amd64) and Apple Silicon (arm64)
- **Linux**: Intel/AMD (amd64) and ARM (arm64)  
- **Windows**: Intel/AMD (amd64) and ARM (arm64)

### Installation Methods

#### 1. One-Line Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh
```

#### 2. Manual Download
- Download appropriate binary from `bin/` directory
- Copy to `~/.local/bin/docs4context-com`
- Make executable and add to PATH

#### 3. Build from Source
```bash
git clone https://github.com/jasonwillschiu/docs4context-com
cd docs4context-com
go build -o docs4context-com .
```

### MCP Client Configuration

#### Claude Code
Use the CLI command to add the server:
```bash
# For local installation
claude mcp add docs4context -- docs4context-com

# Alternative: Remote server (if available)
# claude mcp add --transport sse docs4context https://your-server-url/sse
```

#### Cursor
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

Alternative configurations:
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

#### Windsurf
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

#### OpenCode
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

## MCP Tools Available

### Document Management
- **`save_context_document`** - Downloads pre-processed llms.txt files from context7.com with accurate token counting

### Search & Discovery Tools
- **`search_titles`** - Find topics by title keywords with optional repo filtering in local documents
- **`search_content`** - Search across descriptions and code content in locally stored context documents
- **`get_topic_details`** - Extract complete topic information from specific line numbers in local files
- **`list_repositories`** - Show all locally downloaded repositories with metadata and topic counts
- **`analyze_keywords`** - Cross-repository keyword frequency analysis across local documents with relevance scoring

## Context Document Flow

1. Tool receives GitHub URL (accepts both full URLs and `username/repo` format)
2. Downloads pre-processed context document from `https://context7.com/[username]/[repo]/llms.txt`
3. **Counts actual tokens** using tiktoken-go with cl100k_base encoding (GPT-4 compatible)
4. **Saves locally with metadata header**:
   ```
   # METADATA
   # TOKEN_COUNT: 66551
   # DATE_CREATED: 2025-06-26T10:30:45Z
   # REPO: mark3labs/mcp-go
   # SOURCE: https://context7.com/mark3labs/mcp-go/llms.txt
   #
   ```
5. Stores to local `llm-context/[username]/[repo]/llms.txt` with proper permissions
6. Provides basic search and filtering capabilities across locally stored documents

### MCP Protocol Requirements
- All logging must go to stderr (stdout interferes with MCP communication)
- Tool responses must include operation results and file paths
- Error handling should be comprehensive with descriptive messages
- Search tools automatically skip metadata header lines when parsing local content
- All operations work on locally stored context documents

### Configuration
The project integrates with MCP clients via standard configuration:
- `docs4context` points to main server binary
- Can be used alongside other MCP servers
- Supports both local and remote MCP configurations

Example configurations for different clients:

**Claude Code**:
```bash
claude mcp add docs4context -- docs4context-com
```

**Cursor (`mcp.json`)**:
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

**Windsurf**:
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

**OpenCode (`opencode.json`)**:
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

## Usage Patterns for AI Agents

### Repository Discovery Workflow
```
1. Use list_repositories to see locally available context documents
2. Use analyze_keywords to find local repos relevant to user query
3. Use search_titles for specific topic discovery in local documents
4. Use search_content for detailed code/implementation search in local files
5. Use get_topic_details to extract complete information from local context
```

### Example Agent Commands
```
# Discover what's available
→ use list_repositories

# Find repositories about authentication
→ use analyze_keywords with keyword="auth"

# Search for specific implementation patterns
→ use search_titles with query="middleware" repo_filter="mark3labs/mcp-go"

# Get detailed information about specific findings
→ use get_topic_details with repo="mark3labs/mcp-go" line_numbers="45,123,200"
```

### Advanced Usage Patterns

#### Learning New Framework Pattern
```
1. save_context_document("framework/repo") - downloads llms.txt from context7.com
2. search_titles("getting started") - searches local document
3. search_content("example implementation") - searches local content
4. get_topic_details for specific implementation examples from local file
5. search_titles("best practices") for advanced patterns in local document
```

#### Cross-Repository Analysis Pattern  
```
1. save_context_document for multiple related repositories (downloads multiple llms.txt files)
2. analyze_keywords("error handling") to find relevant local documents
3. search_content("error patterns") across all local context files
4. get_topic_details for specific implementations from local documents
5. Compare approaches across different locally stored repositories
```

#### API Discovery Pattern
```
1. list_repositories to see locally available options
2. search_titles("API") or search_titles("endpoints") in local documents
3. search_content("REST API setup") for implementation details in local files
4. get_topic_details for complete API examples from local context
5. search_content("authentication") for security patterns in local documents
```

### Real-World Usage Examples

#### Example 1: Learning Go Web Development
```
User: "I want to learn how to build REST APIs in Go"

AI Agent workflow:
1. save_context_document("gin-gonic/gin")
2. save_context_document("gorilla/mux") 
3. save_context_document("labstack/echo")
4. search_titles("REST API") across all repos
5. search_content("route handler") for implementation patterns
6. get_topic_details for specific API examples
7. analyze_keywords("middleware") to understand middleware patterns
```

#### Example 2: Understanding Authentication Patterns
```
User: "Show me different authentication approaches in Go"

AI Agent workflow:
1. analyze_keywords("authentication") to find relevant repos
2. search_titles("auth") across identified repositories
3. search_content("JWT token") for token-based auth
4. search_content("session") for session-based auth
5. get_topic_details for complete authentication implementations
6. Compare different approaches across repositories
```

#### Example 3: Database Integration Research
```
User: "How do different projects handle database connections?"

AI Agent workflow:
1. analyze_keywords("database") to identify relevant repos
2. search_titles("database connection") for setup patterns
3. search_content("connection pool") for connection management
4. search_content("migration") for database migration patterns
5. get_topic_details for complete database setup examples
6. search_content("ORM") vs search_content("raw SQL") for different approaches
```

## Advantages Over Bash Functions

**Universal Compatibility:** Search tools work in any project via MCP - no need to copy bash functions to other repositories' CLAUDE.md files

**Agent-Friendly:** AI agents can directly call MCP tools instead of executing bash commands

**Rich Metadata:** Each locally stored repository includes accurate token counts, download timestamps, and source information

**Context-Aware:** Search results include line numbers and surrounding context for precise code location in local files

**Cross-Project:** Single MCP server can search across all locally downloaded repositories regardless of current working directory

**Local Access:** All operations work on locally stored context documents, no network calls needed for search operations

## Versioning & Release Management

### Version Information
- **Version Source**: `changelog.md` is the single source of truth
- **Version Format**: Semantic versioning (e.g., `0.1.2`)
- **Build Info**: Embedded in binaries via Go build flags
- **Auto-Update**: Built-in update mechanism via GitHub releases

### Development Workflow
```bash
# Local development
bun run cicd.js --mode dev

# Build for current platform
bun run cicd.js --mode build

# Cross-platform builds with version info
bun run cicd.js --build

# Git operations
bun run cicd.js --commit --tag --push

# Create GitHub release
bun run cicd.js --release

# Full release workflow
bun run cicd.js --build --commit --tag --push --release
```

### Release Process
1. **Update Version**: Edit `changelog.md` with new version and changes
2. **Build & Release**: Run `bun run cicd.js --build --commit --tag --push --release`
3. **Verify**: Check GitHub release page and test auto-update
4. **Announce**: Update documentation and notify users

### Version Commands
```bash
# Check version information
docs4context-com --version

# Check for available updates
docs4context-com --check-updates

# Update to latest version
docs4context-com --update

# Show help
docs4context-com --help
```

### Troubleshooting Updates
- **Update fails**: Binary creates backup automatically, restores on failure
- **Version mismatch**: Use `--check-updates` to verify latest version
- **Network issues**: Install script supports both curl and wget
- **Permission errors**: Ensure binary has execute permissions
- **GitHub API limits**: Rate limits may affect update checks

## Migration from Bash Functions

The previous bash functions in this file have been replaced by MCP tools that provide superior functionality:

- `get_repo_profiles` → `list_repositories` tool
- `analyze_keyword_across_repos` → `analyze_keywords` tool  
- `search_with_context` → `search_titles` + `search_content` tools
- `extract_topics_multi_repo` → `get_topic_details` tool

AI agents should use the MCP tools instead of bash functions for better integration and consistency across projects.