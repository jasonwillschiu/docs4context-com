# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a collection of MCP (Model Context Protocol) servers providing tools for LLM agents:

1. **Main docs4context server** (`main.go`) - Downloads repository context documents from GitHub via context7.com with token counting and semantic search capabilities
2. **Calculator MCP server** (`calculator-mcp/main.go`) - Provides arithmetic operations

Both servers use the `github.com/mark3labs/mcp-go` framework with stdio transport.

## Development Commands

### Build
```bash
# Build main server
go build -o docs4context-com .

# Build calculator server  
cd calculator-mcp && go build -o calculator-mcp .

# Build both
go build -o docs4context-com . && cd calculator-mcp && go build -o calculator-mcp .
```

### Test
```bash
# Run calculator integration test
./test-calculator-integration.sh

# Install/update dependencies
go mod tidy
```

### Run
```bash
# Run main docs4context server
./docs4context-com

# Run calculator server
./calculator-mcp/calculator-mcp
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

## MCP Tools Available

### Document Management
- **`save_context_document`** - Downloads repository context with accurate token counting

### Search & Discovery Tools
- **`search_titles`** - Find topics by title keywords with optional repo filtering
- **`search_content`** - Search across descriptions and code content with context
- **`get_topic_details`** - Extract complete topic information from specific line numbers
- **`list_repositories`** - Show all available repositories with metadata and topic counts
- **`analyze_keywords`** - Cross-repository keyword frequency analysis with relevance scoring

## Enhanced Context Document Flow

1. Tool receives GitHub URL (accepts both full URLs and `username/repo` format)
2. Downloads context document from `https://context7.com/[username]/[repo]/llms.txt`
3. **Counts actual tokens** using tiktoken-go with cl100k_base encoding (GPT-4 compatible)
4. **Saves with metadata header**:
   ```
   # METADATA
   # TOKEN_COUNT: 66551
   # DATE_CREATED: 2025-06-26T10:30:45Z
   # REPO: mark3labs/mcp-go
   # SOURCE: https://context7.com/mark3labs/mcp-go/llms.txt
   #
   ```
5. Stores to `llm-context/[username]/[repo]/llms.txt` with proper permissions

### MCP Protocol Requirements
- All logging must go to stderr (stdout interferes with MCP communication)
- Tool responses must include operation results and file paths
- Error handling should be comprehensive with descriptive messages
- Search tools automatically skip metadata header lines when parsing content

### Configuration
The project integrates with opencode via `opencode.json`:
- `document-downloader` points to main server
- `localmcp` points to calculator server  
- `remotemcp` connects to context7.com MCP service

## Usage Patterns for AI Agents

### Repository Discovery Workflow
```
1. Use list_repositories to see available context documents
2. Use analyze_keywords to find repos relevant to user query
3. Use search_titles for specific topic discovery
4. Use search_content for detailed code/implementation search
5. Use get_topic_details to extract complete information
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

## Advantages Over Bash Functions

**Universal Compatibility:** Search tools work in any project via MCP - no need to copy bash functions to other repositories' CLAUDE.md files

**Agent-Friendly:** AI agents can directly call MCP tools instead of executing bash commands

**Rich Metadata:** Each repository includes accurate token counts, download timestamps, and semantic analysis

**Context-Aware:** Search results include line numbers and surrounding context for precise code location

**Cross-Project:** Single MCP server can search across all downloaded repositories regardless of current working directory

## Migration from Bash Functions

The previous bash functions in this file have been replaced by MCP tools that provide superior functionality:

- `get_repo_profiles` → `list_repositories` tool
- `analyze_keyword_across_repos` → `analyze_keywords` tool  
- `search_with_context` → `search_titles` + `search_content` tools
- `extract_topics_multi_repo` → `get_topic_details` tool

AI agents should use the MCP tools instead of bash functions for better integration and consistency across projects.