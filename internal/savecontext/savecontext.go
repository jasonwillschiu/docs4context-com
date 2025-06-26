package savecontext

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkoukk/tiktoken-go"
)

const defaultTokenCount = 100000000 // 100 million tokens

// AddTool adds the document saving tool to the server
func AddTool(s *server.MCPServer) {
	saveContextTool := mcp.NewTool("save_context_document",
		mcp.WithDescription("Download and save repository context document for LLM use from GitHub URL"),
		mcp.WithString("github_url",
			mcp.Required(),
			mcp.Description("GitHub repository URL (e.g., https://github.com/nanostores/nanostores or nanostores/nanostores)"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for saving the context document (defaults to 'llm-context')"),
		),
	)

	s.AddTool(saveContextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("SAVE_CONTEXT_DOCUMENT tool called")

		githubURL, err := request.RequireString("github_url")
		if err != nil {
			log.Printf("SAVE_CONTEXT_DOCUMENT tool error - invalid parameter 'github_url': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		outputDir := "llm-context"
		if outputDirParam := request.GetString("output_dir", ""); outputDirParam != "" {
			outputDir = outputDirParam
		}

		// Parse GitHub URL to extract username and repository
		username, repo, err := ParseGitHubURL(githubURL)
		if err != nil {
			log.Printf("SAVE_CONTEXT_DOCUMENT tool error - failed to parse GitHub URL: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse GitHub URL: %v", err)), nil
		}

		log.Printf("Parsed GitHub URL: username=%s, repo=%s", username, repo)

		// Try to download the document directly first (bypassing token count for now)
		var tokenCount int
		log.Printf("Attempting direct download without token count...")
		content, err := downloadContextDocument(username, repo, 0)
		if err != nil {
			log.Printf("Direct download failed: %v. Trying with token count...", err)
			// If direct download fails, try to get token count first
			tokenCount, tokenErr := fetchTokenCount(username, repo)
			if tokenErr != nil {
			log.Printf("SAVE_CONTEXT_DOCUMENT tool - failed to fetch exact token count, using %d tokens as fallback: %v (original error: %v)", defaultTokenCount, tokenErr, err)
			tokenCount = defaultTokenCount			}

			log.Printf("Token count retrieved: %d", tokenCount)

			// Download the context document with token count
			content, err = downloadContextDocument(username, repo, tokenCount)
			if err != nil {
				log.Printf("SAVE_CONTEXT_DOCUMENT tool error - failed to download document: %v", err)
				return mcp.NewToolResultError(fmt.Sprintf("failed to download document: %v", err)), nil
			}
		} else {
			log.Printf("Direct download successful!")
		}

		// Count tokens using tiktoken
		actualTokenCount, err := countTokens(content)
		if err != nil {
			log.Printf("SAVE_CONTEXT_DOCUMENT tool warning - failed to count tokens: %v, using default count", err)
			actualTokenCount = tokenCount // Use the original count as fallback
		}

		// Save the document to the specified directory with metadata
		outputPath := filepath.Join(outputDir, username, repo, "llms.txt")
		err = saveDocument(outputPath, content, username, repo, actualTokenCount)
		if err != nil {
			log.Printf("SAVE_CONTEXT_DOCUMENT tool error - failed to save document: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("failed to save document: %v", err)), nil
		}

		log.Printf("SAVE_CONTEXT_DOCUMENT tool: Successfully saved context document to %s (%d tokens)", outputPath, actualTokenCount)
		return mcp.NewToolResultText(fmt.Sprintf("Successfully downloaded and saved context document to %s\nTokens: %d\nSize: %d bytes", outputPath, actualTokenCount, len(content))), nil
	})
}

// ParseGitHubURL extracts username and repository name from various GitHub URL formats
func ParseGitHubURL(url string) (username, repo string, err error) {
	// Remove any trailing slashes
	url = strings.TrimSuffix(url, "/")

	// Handle different URL formats
	if strings.HasPrefix(url, "https://github.com/") {
		// Full GitHub URL: https://github.com/username/repo
		parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], nil
		}
	} else if strings.Count(url, "/") == 1 && !strings.Contains(url, "://") {
		// Short format: username/repo
		parts := strings.Split(url, "/")
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}

	return "", "", fmt.Errorf("invalid GitHub URL format. Expected: https://github.com/username/repo or username/repo")
}

// countTokens counts the number of tokens in the given content using tiktoken
func countTokens(content []byte) (int, error) {
	// Use cl100k_base encoding (used by GPT-4, GPT-3.5-turbo)
	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return 0, fmt.Errorf("failed to get encoding: %v", err)
	}

	// Count tokens
	tokens := encoding.Encode(string(content), nil, nil)
	return len(tokens), nil
}

// fetchTokenCount retrieves the token count from context7.com
func fetchTokenCount(username, repo string) (int, error) {
	url := fmt.Sprintf("https://context7.com/%s/%s", username, repo)
	log.Printf("Fetching token count from: %s", url)

	// Create a client with a User-Agent header to avoid being blocked
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch page, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	// Extract token count using multiple regex patterns
	// First try: HTML structure with spans containing "Tokens:" followed by the number
	// Pattern matches: <span>Tokens:</span><span>66,551</span>
	tokenRegex := regexp.MustCompile(`<span[^>]*>Tokens:</span><span[^>]*>([0-9,]+)</span>`)
	matches := tokenRegex.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		// Second try: Look for "Tokens:" with optional whitespace and number in JSON or text
		tokenRegex = regexp.MustCompile(`"?Tokens"?\s*:?\s*"?([0-9,]+)"?`)
		matches = tokenRegex.FindStringSubmatch(string(body))
	}

	if len(matches) < 2 {
		// Third try: More flexible pattern for Tokens followed by number
		tokenRegex = regexp.MustCompile(`(?i)tokens[:\s]*([0-9,]+)`)
		matches = tokenRegex.FindStringSubmatch(string(body))
	}

	if len(matches) < 2 {
		// Final fallback: Look for any sequence that might contain token info
		preview := string(body)
		if len(preview) > 1000 {
			preview = preview[:1000]
		}
		log.Printf("Debug: Page content preview: %s", preview)
		log.Printf("Token count not found in page content, will use %d as fallback", defaultTokenCount)
		return defaultTokenCount, nil
	}

	// Remove commas from the number and convert to int
	tokenStr := strings.ReplaceAll(matches[1], ",", "")
	tokenCount, err := strconv.Atoi(tokenStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse token count: %v", err)
	}

	return tokenCount, nil
}

// downloadContextDocument downloads the llms.txt file with the specified token count
func downloadContextDocument(username, repo string, tokenCount int) ([]byte, error) {
	var url string
	if tokenCount > 0 {
		url = fmt.Sprintf("https://context7.com/%s/%s/llms.txt?tokens=%d", username, repo, tokenCount)
	} else {
		url = fmt.Sprintf("https://context7.com/%s/%s/llms.txt?tokens=%d", username, repo, defaultTokenCount)
	}
	log.Printf("Downloading context document from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download document, status: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read document content: %v", err)
	}

	return content, nil
}

// saveDocument saves the downloaded content to the specified path with metadata header
func saveDocument(outputPath string, content []byte, username, repo string, tokenCount int) error {
	// Create the directory structure if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Generate metadata header
	currentTime := time.Now().UTC().Format(time.RFC3339)
	sourceURL := fmt.Sprintf("https://context7.com/%s/%s/llms.txt", username, repo)
	
	header := fmt.Sprintf(`# METADATA
# TOKEN_COUNT: %d
# DATE_CREATED: %s
# REPO: %s/%s
# SOURCE: %s
#
`, tokenCount, currentTime, username, repo, sourceURL)

	// Combine header with content
	finalContent := []byte(header)
	finalContent = append(finalContent, content...)

	// Write the content to the file
	if err := os.WriteFile(outputPath, finalContent, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", outputPath, err)
	}

	return nil
}