#!/bin/bash
set -e

# docs4context MCP Server Install Script
# Downloads and installs the appropriate binary for your platform

REPO_URL="https://raw.githubusercontent.com/jasonwillschiu/docs4context-com/main"
GITHUB_API_URL="https://api.github.com/repos/jasonwillschiu/docs4context-com/releases"
INSTALL_DIR="$(pwd)"
BINARY_NAME="docs4context-com"
VERSION=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Darwin*) os="darwin" ;;
        Linux*)  os="linux" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version from GitHub API
get_latest_version() {
    local latest_url="${GITHUB_API_URL}/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -s "$latest_url" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$latest_url" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Check if binary exists and get current version
get_current_version() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        binary_path="${binary_path}.exe"
    fi
    
    if [[ -f "$binary_path" ]]; then
        if "$binary_path" --version 2>/dev/null | head -n1 | cut -d' ' -f2; then
            return 0
        fi
    fi
    
    echo ""
}

# Download and install binary
install_binary() {
    local platform="$1"
    local version="$2"
    local binary_suffix=""
    
    if [[ "$platform" == *"windows"* ]]; then
        binary_suffix=".exe"
    fi
    
    local remote_binary="${BINARY_NAME}-${platform}${binary_suffix}"
    local local_binary="${INSTALL_DIR}/${BINARY_NAME}${binary_suffix}"
    
    # Determine download URL based on version
    local download_url
    if [[ -n "$version" ]]; then
        download_url="https://github.com/jasonwillschiu/docs4context-com/releases/download/v${version}/${remote_binary}"
        log_info "Installing version: $version"
    else
        download_url="${REPO_URL}/bin/${remote_binary}"
        log_info "Installing latest development version"
    fi
    
    log_info "Detected platform: $platform"
    log_info "Downloading from: $download_url"
    
    # Create backup if existing binary exists
    if [[ -f "$local_binary" ]]; then
        log_info "Creating backup of existing binary..."
        cp "$local_binary" "${local_binary}.backup"
    fi
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        if ! curl -fsSL "$download_url" -o "$local_binary"; then
            log_error "Failed to download binary from: $download_url"
            # Restore backup if it exists
            if [[ -f "${local_binary}.backup" ]]; then
                mv "${local_binary}.backup" "$local_binary"
                log_info "Restored backup binary"
            fi
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -qO "$local_binary" "$download_url"; then
            log_error "Failed to download binary from: $download_url"
            # Restore backup if it exists
            if [[ -f "${local_binary}.backup" ]]; then
                mv "${local_binary}.backup" "$local_binary"
                log_info "Restored backup binary"
            fi
            exit 1
        fi
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Make executable
    chmod +x "$local_binary"
    
    # Remove backup on successful install
    if [[ -f "${local_binary}.backup" ]]; then
        rm "${local_binary}.backup"
    fi
    
    log_success "Binary installed to: $local_binary"
}

# Check if install directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_warning "Install directory $INSTALL_DIR is not in your PATH"
        log_info "Add this line to your shell profile (.bashrc, .zshrc, etc.):"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
}

# Show MCP configuration example
show_mcp_config() {
    log_info "MCP Client Configuration Example:"
    echo ""
    echo "Add this to your MCP client config (e.g., opencode.json):"
    echo ""
    echo "{"
    echo "  \"mcp\": {"
    echo "    \"docs4context\": {"
    echo "      \"type\": \"local\","
    echo "      \"command\": [\"./$BINARY_NAME\"],"
    echo "      \"environment\": {}"
    echo "    }"
    echo "  }"
    echo "}"
    echo ""
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version=*)
                VERSION="${1#*=}"
                shift
                ;;
            --version)
                if [[ -n "$2" ]] && [[ "$2" != --* ]]; then
                    VERSION="$2"
                    shift 2
                else
                    log_error "--version requires a value"
                    exit 1
                fi
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Show help information
show_help() {
    echo "docs4context MCP Server Install Script"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --version <version>  Install specific version (e.g., --version 0.1.2)"
    echo "  --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                   # Install latest release"
    echo "  $0 --version 0.1.2   # Install specific version"
}

# Main installation process
main() {
    parse_args "$@"
    
    local platform
    platform=$(detect_platform)
    
    # Get current version if binary exists
    local current_version
    current_version=$(get_current_version)
    
    # Determine target version
    local target_version="$VERSION"
    if [[ -z "$target_version" ]]; then
        log_info "Fetching latest version from GitHub..."
        target_version=$(get_latest_version)
        if [[ -z "$target_version" ]]; then
            log_warning "Could not determine latest version, installing development version"
        fi
    fi
    
    # Check if upgrade is needed
    if [[ -n "$current_version" ]] && [[ -n "$target_version" ]] && [[ "$current_version" == "$target_version" ]]; then
        log_info "Already running version $current_version"
        log_info "Use --version to install a different version"
        exit 0
    fi
    
    # Show version information
    if [[ -n "$current_version" ]]; then
        log_info "Current version: $current_version"
    fi
    if [[ -n "$target_version" ]]; then
        log_info "Target version: $target_version"
    fi
    
    log_info "Installing docs4context MCP server..."
    
    # Install binary
    if ! install_binary "$platform" "$target_version"; then
        log_error "Installation failed"
        exit 1
    fi
    
    # Verify installation
    local installed_version
    installed_version=$(get_current_version)
    if [[ -n "$installed_version" ]]; then
        log_success "Successfully installed version: $installed_version"
    fi
    
    # Check PATH
    check_path
    
    # Show configuration
    show_mcp_config
    
    log_success "Installation completed successfully!"
    log_info "Run '$BINARY_NAME --version' to verify installation"
    log_info "Run '$BINARY_NAME' to start the MCP server"
}

# Run main function
main "$@"