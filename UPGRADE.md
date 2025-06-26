# Upgrade Guide

This document provides comprehensive guidance for upgrading the docs4context MCP server.

## Quick Upgrade

### Automatic Update (Recommended)
```bash
# Check for available updates
docs4context-com --check-updates

# Update to latest version
docs4context-com --update
```

### Manual Install Script
```bash
# Upgrade to latest release
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh

# Install specific version
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh -s -- --version 0.1.2
```

## Upgrade Methods Comparison

| Method | Pros | Cons | Best For |
|--------|------|------|----------|
| **Auto-update (`--update`)** | ✅ Built-in, fast, atomic<br>✅ Automatic backup/restore<br>✅ Version verification | ❌ Requires GitHub API access<br>❌ Network dependent | Daily users, production |
| **Install script** | ✅ Fresh download<br>✅ Version selection<br>✅ Platform detection | ❌ Requires curl/wget<br>❌ Manual execution | Initial install, specific versions |

## Pre-Upgrade Checklist

### 1. Check Current Version
```bash
docs4context-com --version
```

### 2. Verify MCP Client Configuration
Ensure your MCP client (Claude Code, etc.) points to the correct binary path:
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

### 3. Backup MCP Data (Optional)
```bash
# Backup downloaded context documents
cp -r llm-context llm-context.backup
```

## Version Compatibility

### Breaking Changes
- **v0.1.x → v0.2.x**: No breaking changes expected
- **Major versions**: Check changelog.md for migration notes

### MCP Client Compatibility
- **Binary name**: Always `docs4context-com` (no version suffix)
- **Command interface**: `--version`, `--help`, `--update` flags maintained
- **MCP protocol**: Compatible across versions

## Post-Upgrade Verification

### 1. Version Check
```bash
docs4context-com --version
```
Expected output:
```
docs4context-com 0.1.2
Build Date: 2025-06-26T10:30:45Z
Git Commit: abc1234
Go Version: go1.24.2
OS/Arch: darwin/arm64
```

### 2. MCP Server Test
```bash
# Test MCP server starts correctly
docs4context-com
# Should show: "Starting docs4context MCP Server 0.1.2"
# Press Ctrl+C to stop
```

### 3. Client Integration Test
Restart your MCP client and verify:
- docs4context tools are available
- Existing context documents still accessible
- New features work as expected

## Troubleshooting

### Update Command Fails

**Error**: `Error checking for updates: failed to fetch latest release`
```bash
# Solutions:
1. Check internet connection
2. Verify GitHub API access
3. Check for rate limiting
4. Use manual install method as fallback
```

**Error**: `Error updating: no binary found for platform`
```bash
# Solutions:
1. Check platform support: darwin/linux/windows
2. Verify architecture: amd64/arm64
3. Use manual download from GitHub releases
```

**Error**: `Error updating: failed to replace binary`
```bash
# Solutions:
1. Check file permissions: chmod +x docs4context-com
2. Ensure binary not in use (stop MCP clients)
3. Run with elevated privileges if needed
4. Manual install to different location
```

### Install Script Fails

**Error**: `Failed to download binary`
```bash
# Solutions:
1. Check download URL accessibility
2. Verify version exists in GitHub releases
3. Check curl/wget installation
4. Try alternative download method
```

**Error**: `Neither curl nor wget found`
```bash
# Install required tools:
# macOS: brew install curl
# Linux: sudo apt-get install curl
# Windows: Install via package manager
```

### Version Mismatch

**Issue**: Binary shows old version after update
```bash
# Solutions:
1. Verify correct binary path in use
2. Check for multiple installations
3. Clear shell cache: hash -r
4. Restart terminal session
```

### Permission Errors

**Error**: `Permission denied` when running binary
```bash
# Solutions:
chmod +x docs4context-com
# Or reinstall with correct permissions
```

### GitHub API Rate Limiting

**Error**: `API rate limit exceeded`
```bash
# Solutions:
1. Wait for rate limit reset (usually 1 hour)
2. Use manual install method
3. Authenticate with GitHub CLI: gh auth login
```

## Rollback Procedure

### Automatic Backup
The auto-updater creates automatic backups:
```bash
# If update fails, backup is automatically restored
# No manual action needed
```

### Manual Rollback
```bash
# Install specific previous version
curl -fsSL https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main/install.sh | sh -s -- --version 0.1.1

# Or download manually from GitHub releases
# https://github.com/jasonwillschiu/docs4context-com/releases
```

## Advanced Upgrade Scenarios

### Multiple Installation Locations
```bash
# Check all installations
which -a docs4context-com

# Update specific installation
/path/to/docs4context-com --update
```

### Corporate Networks
```bash
# If direct GitHub access blocked:
1. Download binary manually from releases page
2. Transfer to target system
3. Replace existing binary
4. Verify with --version
```

### Automation
```bash
#!/bin/bash
# Automated upgrade script
set -e

echo "Checking for updates..."
if docs4context-com --check-updates | grep -q "New version available"; then
    echo "Updating docs4context-com..."
    docs4context-com --update
    echo "Update completed"
else
    echo "Already up to date"
fi
```

## Best Practices

### 1. Regular Updates
- Check for updates weekly: `docs4context-com --check-updates`
- Subscribe to GitHub releases for notifications
- Read changelog.md before major updates

### 2. Staging Updates
- Test updates in development environment first
- Verify MCP client compatibility
- Check for breaking changes in release notes

### 3. Backup Strategy
- Keep backup of working binary for critical systems
- Document current configuration before updates
- Test rollback procedures periodically

### 4. Monitoring
- Monitor MCP server logs after updates
- Verify all expected tools are available
- Check performance impact of new versions

## Getting Help

### Support Channels
- **GitHub Issues**: https://github.com/jasonwillschiu/docs4context-com/issues
- **Documentation**: https://github.com/jasonwillschiu/docs4context-com/blob/main/CLAUDE.md
- **Releases**: https://github.com/jasonwillschiu/docs4context-com/releases

### Reporting Issues
When reporting upgrade issues, include:
- Current version: `docs4context-com --version`
- Target version
- Full error message
- Operating system and architecture
- MCP client type and version

### Emergency Rollback
If critical issues occur after upgrade:
1. Stop all MCP clients
2. Restore previous version immediately
3. Restart MCP clients
4. Report issue on GitHub
5. Investigate in non-production environment