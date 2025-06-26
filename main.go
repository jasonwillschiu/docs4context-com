package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"docs4context-com/internal/savecontext"
	"docs4context-com/internal/search"
	"docs4context-com/internal/updater"

	"github.com/mark3labs/mcp-go/server"
)

// Build information (set via ldflags during build)
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Command line flags
	var (
		showVersion   = flag.Bool("version", false, "Show version information")
		showHelp      = flag.Bool("help", false, "Show help information")
		updateBinary  = flag.Bool("update", false, "Check for and install updates")
		checkUpdates  = flag.Bool("check-updates", false, "Check for available updates without installing")
	)
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("docs4context-com %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	// Handle help flag
	if *showHelp {
		fmt.Println("docs4context MCP Server")
		fmt.Println("Usage:")
		fmt.Println("  docs4context-com [options]")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --version         Show version information")
		fmt.Println("  --help            Show this help message")
		fmt.Println("  --update          Check for and install updates")
		fmt.Println("  --check-updates   Check for available updates without installing")
		fmt.Println("")
		fmt.Println("This is an MCP (Model Context Protocol) server that provides")
		fmt.Println("document context and search tools for AI agents.")
		return
	}

	// Handle check updates flag
	if *checkUpdates {
		fmt.Println("Checking for updates...")
		release, hasUpdate, err := updater.CheckForUpdates(Version)
		if err != nil {
			fmt.Printf("Error checking for updates: %v\n", err)
			os.Exit(1)
		}
		
		if hasUpdate {
			fmt.Printf("New version available: %s\n", release.TagName)
			fmt.Printf("Current version: %s\n", Version)
			fmt.Println("Run with --update to install the latest version")
		} else {
			fmt.Println("You are running the latest version")
		}
		return
	}

	// Handle update flag
	if *updateBinary {
		fmt.Println("Checking for updates...")
		release, hasUpdate, err := updater.CheckForUpdates(Version)
		if err != nil {
			fmt.Printf("Error checking for updates: %v\n", err)
			os.Exit(1)
		}
		
		if !hasUpdate {
			fmt.Println("You are already running the latest version")
			return
		}
		
		fmt.Printf("Updating from %s to %s...\n", Version, release.TagName)
		if err := updater.DownloadUpdate(release); err != nil {
			fmt.Printf("Error updating: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Update completed successfully!")
		fmt.Println("Please restart the application to use the new version")
		return
	}

	// Set up logging to stderr so it doesn't interfere with stdio communication
	log.SetOutput(os.Stderr)
	log.SetPrefix("[docs4context] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Starting docs4context MCP Server %s", Version)

	// Create a new MCP server
	s := server.NewMCPServer(
		"docs4context",
		Version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add the document context saving tool
	log.Println("Registering save_context_document tool...")
	savecontext.AddTool(s)
	log.Println("Tool registered successfully")

	// Add search tools
	log.Println("Registering search tools...")
	search.AddSearchTitles(s)
	search.AddSearchContent(s)
	search.AddGetTopicDetails(s)
	search.AddListRepositories(s)
	search.AddAnalyzeKeywords(s)
	log.Println("Search tools registered successfully")

	// Start the stdio server
	log.Println("Starting stdio server...")
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
		fmt.Printf("Server error: %v\n", err)
	}
}
