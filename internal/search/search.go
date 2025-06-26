package search

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSearchTitles adds the search titles tool to the server
func AddSearchTitles(s *server.MCPServer) {
	searchTool := mcp.NewTool("search_titles",
		mcp.WithDescription("Search for topics by title keywords across downloaded repository context documents"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to match against topic titles"),
		),
		mcp.WithString("repo_filter",
			mcp.Description("Optional repository filter in format 'username/repo' to limit search scope"),
		),
	)

	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("SEARCH_TITLES tool called")

		query, err := request.RequireString("query")
		if err != nil {
			log.Printf("SEARCH_TITLES tool error - invalid parameter 'query': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		repoFilter := request.GetString("repo_filter", "")

		results, err := searchTitles(query, repoFilter)
		if err != nil {
			log.Printf("SEARCH_TITLES tool error - search failed: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
		}

		log.Printf("SEARCH_TITLES tool: Found %d results for query '%s'", len(results), query)
		return mcp.NewToolResultText(results), nil
	})
}

// searchTitles searches for topics by title keywords
func searchTitles(query, repoFilter string) (string, error) {
	contextDir := "llm-context"
	
	// Check if context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return "No context documents found. Please download some repositories first using save_context_document.", nil
	}

	var results []string
	results = append(results, fmt.Sprintf("=== Search Results for Title Query: '%s' ===\n", query))

	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process llms.txt files
		if info.Name() != "llms.txt" {
			return nil
		}

		// Extract repo info from path
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		if len(pathParts) < 3 {
			return nil
		}
		repoName := pathParts[0] + "/" + pathParts[1]

		// Apply repository filter if specified
		if repoFilter != "" && repoName != repoFilter {
			return nil
		}

		// Search within this file
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return nil
		}

		lines := strings.Split(string(content), "\n")
		foundMatches := false

		for i, line := range lines {
			// Skip metadata header lines
			if strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "TITLE:") && strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
				if !foundMatches {
					results = append(results, fmt.Sprintf("\n--- %s ---", repoName))
					foundMatches = true
				}
				results = append(results, fmt.Sprintf("Line %d: %s", i+1, line))
				// Also include the next line if it's a description
				if i+1 < len(lines) && strings.HasPrefix(lines[i+1], "DESCRIPTION:") {
					results = append(results, fmt.Sprintf("Line %d: %s", i+2, lines[i+1]))
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to search files: %v", err)
	}

	if len(results) == 1 {
		results = append(results, "\nNo matching titles found.")
	}

	return strings.Join(results, "\n"), nil
}

// AddSearchContent adds the search content tool to the server
func AddSearchContent(s *server.MCPServer) {
	searchTool := mcp.NewTool("search_content",
		mcp.WithDescription("Search across descriptions and code content in downloaded repository context documents"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to match against descriptions and code content"),
		),
		mcp.WithString("repo_filter",
			mcp.Description("Optional repository filter in format 'username/repo' to limit search scope"),
		),
	)

	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("SEARCH_CONTENT tool called")

		query, err := request.RequireString("query")
		if err != nil {
			log.Printf("SEARCH_CONTENT tool error - invalid parameter 'query': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		repoFilter := request.GetString("repo_filter", "")

		results, err := searchContent(query, repoFilter)
		if err != nil {
			log.Printf("SEARCH_CONTENT tool error - search failed: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
		}

		log.Printf("SEARCH_CONTENT tool: Found results for query '%s'", query)
		return mcp.NewToolResultText(results), nil
	})
}

// searchContent searches across descriptions and code content
func searchContent(query, repoFilter string) (string, error) {
	contextDir := "llm-context"
	
	// Check if context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return "No context documents found. Please download some repositories first using save_context_document.", nil
	}

	var results []string
	results = append(results, fmt.Sprintf("=== Search Results for Content Query: '%s' ===\n", query))

	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process llms.txt files
		if info.Name() != "llms.txt" {
			return nil
		}

		// Extract repo info from path
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		if len(pathParts) < 3 {
			return nil
		}
		repoName := pathParts[0] + "/" + pathParts[1]

		// Apply repository filter if specified
		if repoFilter != "" && repoName != repoFilter {
			return nil
		}

		// Search within this file
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return nil
		}

		lines := strings.Split(string(content), "\n")
		foundMatches := false
		queryLower := strings.ToLower(query)

		for i, line := range lines {
			// Skip metadata header lines
			if strings.HasPrefix(line, "#") {
				continue
			}

			lineLower := strings.ToLower(line)
			
			// Search in DESCRIPTION lines and CODE blocks
			if (strings.HasPrefix(line, "DESCRIPTION:") || strings.HasPrefix(line, "CODE:") || 
				(!strings.HasPrefix(line, "TITLE:") && !strings.HasPrefix(line, "SOURCE:") && 
				 !strings.HasPrefix(line, "LANGUAGE:") && !strings.Contains(line, "----------------------------------------"))) &&
				strings.Contains(lineLower, queryLower) {
				
				if !foundMatches {
					results = append(results, fmt.Sprintf("\n--- %s ---", repoName))
					foundMatches = true
				}
				
				// Include some context around the match
				contextStart := i - 2
				if contextStart < 0 {
					contextStart = 0
				}
				contextEnd := i + 2
				if contextEnd >= len(lines) {
					contextEnd = len(lines) - 1
				}
				
				results = append(results, fmt.Sprintf("Match at line %d:", i+1))
				for j := contextStart; j <= contextEnd; j++ {
					prefix := "  "
					if j == i {
						prefix = "* " // Mark the matching line
					}
					results = append(results, fmt.Sprintf("%s%s", prefix, lines[j]))
				}
				results = append(results, "")
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to search files: %v", err)
	}

	if len(results) == 1 {
		results = append(results, "\nNo matching content found.")
	}

	return strings.Join(results, "\n"), nil
}

// AddGetTopicDetails adds the get topic details tool to the server
func AddGetTopicDetails(s *server.MCPServer) {
	detailsTool := mcp.NewTool("get_topic_details",
		mcp.WithDescription("Extract complete topic information with context from specific line numbers in repository documents"),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository in format 'username/repo'"),
		),
		mcp.WithString("line_numbers",
			mcp.Required(),
			mcp.Description("Comma-separated line numbers to extract topics from (e.g., '45,123,200')"),
		),
	)

	s.AddTool(detailsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("GET_TOPIC_DETAILS tool called")

		repo, err := request.RequireString("repo")
		if err != nil {
			log.Printf("GET_TOPIC_DETAILS tool error - invalid parameter 'repo': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		lineNumbersStr, err := request.RequireString("line_numbers")
		if err != nil {
			log.Printf("GET_TOPIC_DETAILS tool error - invalid parameter 'line_numbers': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		results, err := getTopicDetails(repo, lineNumbersStr)
		if err != nil {
			log.Printf("GET_TOPIC_DETAILS tool error - extraction failed: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("extraction failed: %v", err)), nil
		}

		log.Printf("GET_TOPIC_DETAILS tool: Extracted details for repo '%s'", repo)
		return mcp.NewToolResultText(results), nil
	})
}

// getTopicDetails extracts complete topic information from specific line numbers
func getTopicDetails(repo, lineNumbersStr string) (string, error) {
	contextDir := "llm-context"
	filePath := filepath.Join(contextDir, repo, "llms.txt")
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Sprintf("Repository '%s' not found. Please download it first using save_context_document.", repo), nil
	}

	// Parse line numbers
	lineNumberStrings := strings.Split(lineNumbersStr, ",")
	var lineNumbers []int
	for _, lineStr := range lineNumberStrings {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		lineNum, err := strconv.Atoi(lineStr)
		if err != nil {
			return "", fmt.Errorf("invalid line number: %s", lineStr)
		}
		lineNumbers = append(lineNumbers, lineNum)
	}

	if len(lineNumbers) == 0 {
		return "", fmt.Errorf("no valid line numbers provided")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var results []string
	results = append(results, fmt.Sprintf("=== Topic Details for %s ===\n", repo))

	for _, lineNum := range lineNumbers {
		if lineNum < 1 || lineNum > len(lines) {
			results = append(results, fmt.Sprintf("Line %d: OUT OF RANGE (file has %d lines)", lineNum, len(lines)))
			continue
		}

		// Convert to 0-based index
		lineIndex := lineNum - 1
		line := lines[lineIndex]

		// Skip metadata header lines
		if strings.HasPrefix(line, "#") {
			results = append(results, fmt.Sprintf("Line %d: METADATA HEADER - %s", lineNum, line))
			continue
		}

		// Find the complete topic block starting from this line
		if strings.HasPrefix(line, "TITLE:") {
			results = append(results, fmt.Sprintf("\n--- Topic starting at line %d ---", lineNum))
			
			// Extract the complete topic block
			topicLines := []string{fmt.Sprintf("Line %d: %s", lineNum, line)}
			
			// Look for DESCRIPTION, SOURCE, LANGUAGE, and CODE
			currentIndex := lineIndex + 1
			for currentIndex < len(lines) && !strings.Contains(lines[currentIndex], "----------------------------------------") {
				currentLine := lines[currentIndex]
				if strings.HasPrefix(currentLine, "DESCRIPTION:") ||
				   strings.HasPrefix(currentLine, "SOURCE:") ||
				   strings.HasPrefix(currentLine, "LANGUAGE:") ||
				   strings.HasPrefix(currentLine, "CODE:") ||
				   (!strings.HasPrefix(currentLine, "TITLE:") && currentLine != "") {
					topicLines = append(topicLines, fmt.Sprintf("Line %d: %s", currentIndex+1, currentLine))
				}
				currentIndex++
			}
			
			results = append(results, strings.Join(topicLines, "\n"))
		} else {
			// For non-TITLE lines, provide context
			results = append(results, fmt.Sprintf("\n--- Context around line %d ---", lineNum))
			
			contextStart := lineIndex - 3
			if contextStart < 0 {
				contextStart = 0
			}
			contextEnd := lineIndex + 3
			if contextEnd >= len(lines) {
				contextEnd = len(lines) - 1
			}
			
			for i := contextStart; i <= contextEnd; i++ {
				prefix := "  "
				if i == lineIndex {
					prefix = "* " // Mark the requested line
				}
				results = append(results, fmt.Sprintf("%sLine %d: %s", prefix, i+1, lines[i]))
			}
		}
	}

	return strings.Join(results, "\n"), nil
}

// AddListRepositories adds the list repositories tool to the server
func AddListRepositories(s *server.MCPServer) {
	listTool := mcp.NewTool("list_repositories",
		mcp.WithDescription("List all available repositories with their metadata and topic counts"),
	)

	s.AddTool(listTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("LIST_REPOSITORIES tool called")

		results, err := listRepositories()
		if err != nil {
			log.Printf("LIST_REPOSITORIES tool error - listing failed: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("listing failed: %v", err)), nil
		}

		log.Printf("LIST_REPOSITORIES tool: Listed available repositories")
		return mcp.NewToolResultText(results), nil
	})
}

// listRepositories lists all available repositories with metadata
func listRepositories() (string, error) {
	contextDir := "llm-context"
	
	// Check if context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return "No context documents found. Please download some repositories first using save_context_document.", nil
	}

	var results []string
	results = append(results, "=== Available Repositories ===\n")

	type RepoInfo struct {
		Name       string
		TokenCount int
		DateCreated string
		TopicCount int
		Keywords   []string
	}

	var repos []RepoInfo

	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process llms.txt files
		if info.Name() != "llms.txt" {
			return nil
		}

		// Extract repo info from path
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		if len(pathParts) < 3 {
			return nil
		}
		repoName := pathParts[0] + "/" + pathParts[1]

		// Read file to get metadata and count topics
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return nil
		}

		lines := strings.Split(string(content), "\n")
		
		var tokenCount int
		var dateCreated string
		var topicCount int
		var keywords []string

		// Parse metadata header
		for _, line := range lines {
			if strings.HasPrefix(line, "# TOKEN_COUNT:") {
				tokenStr := strings.TrimSpace(strings.TrimPrefix(line, "# TOKEN_COUNT:"))
				if count, err := strconv.Atoi(tokenStr); err == nil {
					tokenCount = count
				}
			} else if strings.HasPrefix(line, "# DATE_CREATED:") {
				dateCreated = strings.TrimSpace(strings.TrimPrefix(line, "# DATE_CREATED:"))
			} else if strings.HasPrefix(line, "TITLE:") {
				topicCount++
			}
		}

		// Extract common keywords from content
		contentStr := strings.ToLower(string(content))
		commonKeywords := []string{"server", "client", "config", "auth", "api", "tool", "mcp", "go", "typescript", "react", "database", "docker"}
		keywordCounts := make(map[string]int)
		
		for _, keyword := range commonKeywords {
			count := strings.Count(contentStr, keyword)
			if count > 2 { // Only include frequently mentioned keywords
				keywordCounts[keyword] = count
			}
		}

		// Sort keywords by frequency and take top 3
		type KeywordCount struct {
			Keyword string
			Count   int
		}
		var keywordList []KeywordCount
		for k, v := range keywordCounts {
			keywordList = append(keywordList, KeywordCount{k, v})
		}
		// Simple sort by count (descending)
		for i := 0; i < len(keywordList)-1; i++ {
			for j := i+1; j < len(keywordList); j++ {
				if keywordList[j].Count > keywordList[i].Count {
					keywordList[i], keywordList[j] = keywordList[j], keywordList[i]
				}
			}
		}
		
		for i := 0; i < len(keywordList) && i < 3; i++ {
			keywords = append(keywords, fmt.Sprintf("%s(%d)", keywordList[i].Keyword, keywordList[i].Count))
		}

		repos = append(repos, RepoInfo{
			Name:        repoName,
			TokenCount:  tokenCount,
			DateCreated: dateCreated,
			TopicCount:  topicCount,
			Keywords:    keywords,
		})

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to scan repositories: %v", err)
	}

	if len(repos) == 0 {
		results = append(results, "No repositories found.")
	} else {
		for _, repo := range repos {
			results = append(results, fmt.Sprintf("📁 %s", repo.Name))
			results = append(results, fmt.Sprintf("   Topics: %d", repo.TopicCount))
			results = append(results, fmt.Sprintf("   Tokens: %d", repo.TokenCount))
			if repo.DateCreated != "" {
				results = append(results, fmt.Sprintf("   Downloaded: %s", repo.DateCreated))
			}
			if len(repo.Keywords) > 0 {
				results = append(results, fmt.Sprintf("   Keywords: %s", strings.Join(repo.Keywords, " ")))
			}
			results = append(results, "")
		}
	}

	return strings.Join(results, "\n"), nil
}

// AddAnalyzeKeywords adds the analyze keywords tool to the server
func AddAnalyzeKeywords(s *server.MCPServer) {
	analyzeTool := mcp.NewTool("analyze_keywords",
		mcp.WithDescription("Analyze keyword frequency across all repositories"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("Keyword to analyze across all repositories"),
		),
	)

	s.AddTool(analyzeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("ANALYZE_KEYWORDS tool called")

		keyword, err := request.RequireString("keyword")
		if err != nil {
			log.Printf("ANALYZE_KEYWORDS tool error - invalid parameter 'keyword': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		results, err := analyzeKeywords(keyword)
		if err != nil {
			log.Printf("ANALYZE_KEYWORDS tool error - analysis failed: %v", err)
			return mcp.NewToolResultError(fmt.Sprintf("analysis failed: %v", err)), nil
		}

		log.Printf("ANALYZE_KEYWORDS tool: Analyzed keyword '%s'", keyword)
		return mcp.NewToolResultText(results), nil
	})
}

// analyzeKeywords analyzes keyword frequency across all repositories
func analyzeKeywords(keyword string) (string, error) {
	contextDir := "llm-context"
	
	// Check if context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return "No context documents found. Please download some repositories first using save_context_document.", nil
	}

	var results []string
	results = append(results, fmt.Sprintf("=== Keyword Analysis: '%s' ===\n", keyword))

	type RepoMatch struct {
		Name       string
		Matches    int
		TitleMatches int
		DescMatches  int
		CodeMatches  int
		TopicCount int
	}

	var repoMatches []RepoMatch
	keywordLower := strings.ToLower(keyword)

	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process llms.txt files
		if info.Name() != "llms.txt" {
			return nil
		}

		// Extract repo info from path
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		if len(pathParts) < 3 {
			return nil
		}
		repoName := pathParts[0] + "/" + pathParts[1]

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return nil
		}

		lines := strings.Split(string(content), "\n")
		
		var titleMatches, descMatches, codeMatches, topicCount int
		inCodeBlock := false

		for _, line := range lines {
			// Skip metadata header lines
			if strings.HasPrefix(line, "#") {
				continue
			}

			lineLower := strings.ToLower(line)
			lineMatches := strings.Count(lineLower, keywordLower)

			if strings.HasPrefix(line, "TITLE:") {
				topicCount++
				titleMatches += lineMatches
			} else if strings.HasPrefix(line, "DESCRIPTION:") {
				descMatches += lineMatches
			} else if strings.HasPrefix(line, "CODE:") {
				inCodeBlock = true
				codeMatches += lineMatches
			} else if strings.Contains(line, "----------------------------------------") {
				inCodeBlock = false
			} else if inCodeBlock {
				codeMatches += lineMatches
			}
		}

		totalMatches := titleMatches + descMatches + codeMatches
		if totalMatches > 0 {
			repoMatches = append(repoMatches, RepoMatch{
				Name:         repoName,
				Matches:      totalMatches,
				TitleMatches: titleMatches,
				DescMatches:  descMatches,
				CodeMatches:  codeMatches,
				TopicCount:   topicCount,
			})
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to analyze keyword: %v", err)
	}

	if len(repoMatches) == 0 {
		results = append(results, "No matches found across any repositories.")
	} else {
		// Sort repositories by total matches (descending)
		for i := 0; i < len(repoMatches)-1; i++ {
			for j := i+1; j < len(repoMatches); j++ {
				if repoMatches[j].Matches > repoMatches[i].Matches {
					repoMatches[i], repoMatches[j] = repoMatches[j], repoMatches[i]
				}
			}
		}

		results = append(results, fmt.Sprintf("Found %d repositories with matches:", len(repoMatches)))
		results = append(results, "")

		for _, repo := range repoMatches {
			results = append(results, fmt.Sprintf("📁 %s: %d total matches", repo.Name, repo.Matches))
			
			breakdown := []string{}
			if repo.TitleMatches > 0 {
				breakdown = append(breakdown, fmt.Sprintf("titles(%d)", repo.TitleMatches))
			}
			if repo.DescMatches > 0 {
				breakdown = append(breakdown, fmt.Sprintf("descriptions(%d)", repo.DescMatches))
			}
			if repo.CodeMatches > 0 {
				breakdown = append(breakdown, fmt.Sprintf("code(%d)", repo.CodeMatches))
			}
			
			if len(breakdown) > 0 {
				results = append(results, fmt.Sprintf("   Breakdown: %s", strings.Join(breakdown, ", ")))
			}
			
			// Calculate relevance percentage
			relevance := 0.0
			if repo.TopicCount > 0 {
				relevance = (float64(repo.Matches) / float64(repo.TopicCount)) * 100
			}
			results = append(results, fmt.Sprintf("   Relevance: %.1f%% (%d matches in %d topics)", relevance, repo.Matches, repo.TopicCount))
			results = append(results, "")
		}

		// Summary statistics
		totalMatches := 0
		totalRepos := len(repoMatches)
		for _, repo := range repoMatches {
			totalMatches += repo.Matches
		}
		
		results = append(results, "--- Summary ---")
		results = append(results, fmt.Sprintf("Total matches: %d across %d repositories", totalMatches, totalRepos))
		if totalRepos > 0 {
			avgMatches := float64(totalMatches) / float64(totalRepos)
			results = append(results, fmt.Sprintf("Average matches per repository: %.1f", avgMatches))
		}
	}

	return strings.Join(results, "\n"), nil
}