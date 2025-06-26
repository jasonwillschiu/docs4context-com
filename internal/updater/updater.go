package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	// "path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	githubReleasesAPI = "https://api.github.com/repos/jasonwillschiu/docs4context-com/releases/latest"
	githubReleaseURL  = "https://github.com/jasonwillschiu/docs4context-com/releases/download"
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// GetLatestRelease fetches the latest release from GitHub
func GetLatestRelease() (*Release, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(githubReleasesAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release JSON: %w", err)
	}

	return &release, nil
}

// GetPlatformBinary returns the binary name for the current platform
func GetPlatformBinary() string {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	binaryName := fmt.Sprintf("docs4context-com-%s", platform)

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	return binaryName
}

// CompareVersions compares two version strings and returns true if latest is newer than current
func CompareVersions(current, latest string) (bool, error) {
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Parse version strings (assumes format: major.minor.patch)
	currentParts, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version format: %s", current)
	}

	latestParts, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf("invalid latest version format: %s", latest)
	}

	// Compare major, minor, patch in order
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true, nil // latest is newer
		} else if latestParts[i] < currentParts[i] {
			return false, nil // current is newer
		}
		// If equal, continue to next part
	}

	return false, nil // versions are equal
}

// parseVersion parses a semantic version string (e.g., "1.2.3") into [major, minor, patch]
func parseVersion(version string) ([]int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("version must have exactly 3 parts (major.minor.patch)")
	}

	result := make([]int, 3)
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("version part '%s' is not a valid number", part)
		}
		result[i] = num
	}

	return result, nil
}

// DownloadUpdate downloads the latest binary for the current platform
func DownloadUpdate(release *Release) error {
	platformBinary := GetPlatformBinary()

	// Find the asset for current platform
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == platformBinary {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for platform %s-%s", runtime.GOOS, runtime.GOARCH)
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Create backup of current binary
	backupPath := currentExe + ".backup"
	if err := copyFile(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Download new binary
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file
	tempPath := currentExe + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Copy downloaded content to temp file
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write downloaded content: %w", err)
	}

	// Make temp file executable
	if err := os.Chmod(tempPath, 0755); err != nil {
		return fmt.Errorf("failed to make temp file executable: %w", err)
	}

	// Replace current binary with new one
	if err := os.Rename(tempPath, currentExe); err != nil {
		// If rename fails, restore backup
		os.Rename(backupPath, currentExe)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Clean up backup
	os.Remove(backupPath)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// CheckForUpdates checks if a newer version is available
func CheckForUpdates(currentVersion string) (*Release, bool, error) {
	release, err := GetLatestRelease()
	if err != nil {
		return nil, false, err
	}

	hasUpdate, err := CompareVersions(currentVersion, release.TagName)
	if err != nil {
		return nil, false, err
	}

	return release, hasUpdate, nil
}
