package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verifyChecksum bool
	verifyOnline   bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	Long: `Display detailed version, build information, and optionally verify binary integrity.

The version command shows:
  - Version number and build metadata
  - Git commit and build date
  - Go version and platform
  - Binary checksum verification (optional)`,
	RunE: runVersion,
	Example: `  # Show version
  uuidkey version
  
  # Verify binary checksum
  uuidkey version --verify
  
  # Verify against online checksums
  uuidkey version --verify-online`,
}

func init() {
	versionCmd.Flags().BoolVar(&verifyChecksum, "verify", false, "Verify binary checksum")
	versionCmd.Flags().BoolVar(&verifyOnline, "verify-online", false, "Verify against online checksums")
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	// Basic version info
	versionInfo := map[string]interface{}{
		"version":   version,
		"commit":    commit,
		"date":      date,
		"built_by":  builtBy,
		"go":        goVersion,
		"platform":  platform,
	}
	
	// Add build info from runtime/debug if available
	if info, ok := debug.ReadBuildInfo(); ok {
		var vcsRevision, vcsTime, vcsModified string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.time":
				vcsTime = setting.Value
			case "vcs.modified":
				vcsModified = setting.Value
			}
		}
		if vcsRevision != "" {
			versionInfo["vcs_revision"] = vcsRevision
		}
		if vcsTime != "" {
			versionInfo["vcs_time"] = vcsTime
		}
		if vcsModified != "" {
			versionInfo["vcs_modified"] = vcsModified == "true"
		}
	}
	
	// Handle checksum verification
	if verifyChecksum || verifyOnline {
		checksum, err := calculateSelfChecksum()
		if err != nil {
			versionInfo["checksum_error"] = err.Error()
		} else {
			versionInfo["checksum"] = checksum
			
			if verifyOnline {
				valid, err := verifyOnlineChecksum(checksum)
				if err != nil {
					versionInfo["online_verify_error"] = err.Error()
				} else {
					versionInfo["checksum_verified"] = valid
				}
			}
		}
	}
	
	// Output
	out := cmd.OutOrStdout()
	
	if jsonOutput {
		output(versionInfo, "")
	} else {
		_, _ = fmt.Fprintf(out, "uuidkey %s", version)
		if date != "unknown" {
			_, _ = fmt.Fprintf(out, " (%s)", date)
		}
		_, _ = fmt.Fprintln(out)
		
		if !quietMode {
			_, _ = fmt.Fprintln(out, "\nBuild Information:")
			_, _ = fmt.Fprintf(out, "  Commit:    %s\n", commit)
			_, _ = fmt.Fprintf(out, "  Built:     %s\n", date)
			_, _ = fmt.Fprintf(out, "  Built by:  %s\n", builtBy)
			_, _ = fmt.Fprintf(out, "  Go:        %s\n", goVersion)
			_, _ = fmt.Fprintf(out, "  Platform:  %s\n", platform)
			
			if verifyChecksum || verifyOnline {
				checksum, err := calculateSelfChecksum()
				if err != nil {
					_, _ = fmt.Fprintf(out, "\n✗ Checksum error: %v\n", err)
				} else {
					_, _ = fmt.Fprintf(out, "\nChecksum:    %s\n", checksum)
					
					if verifyOnline {
						valid, err := verifyOnlineChecksum(checksum)
						if err != nil {
							_, _ = fmt.Fprintf(out, "✗ Online verification error: %v\n", err)
						} else if valid {
							_, _ = fmt.Fprintln(out, "✓ Binary checksum verified")
						} else {
							_, _ = fmt.Fprintln(out, "✗ Binary checksum mismatch")
						}
					}
				}
			}
		}
	}
	
	return nil
}

// osExecutable is a variable to allow mocking in tests
var osExecutable = os.Executable

// calculateSelfChecksum calculates SHA256 of the running binary
func calculateSelfChecksum() (string, error) {
	executable, err := osExecutable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Resolve symlinks
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}
	
	file, err := os.Open(executable)
	if err != nil {
		return "", fmt.Errorf("failed to open binary: %w", err)
	}
	defer func() { _ = file.Close() }()
	
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}
	
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// verifyOnlineChecksum downloads and verifies checksum from GitHub release
func verifyOnlineChecksum(localChecksum string) (bool, error) {
	if version == "dev" {
		return false, fmt.Errorf("cannot verify development version")
	}
	
	// Construct checksums URL
	url := fmt.Sprintf("https://github.com/agentstation/uuidkey/releases/download/%s/checksums.txt", version)
	
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to download checksums: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("checksums not found (HTTP %d)", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read checksums: %w", err)
	}
	
	// Parse checksums file
	binaryName := fmt.Sprintf("uuidkey_%s_%s", strings.Replace(platform, "/", "_", 1), version)
	if strings.Contains(platform, "windows") {
		binaryName += ".exe"
	}
	
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 && strings.Contains(parts[1], binaryName) {
			return parts[0] == localChecksum, nil
		}
	}
	
	return false, fmt.Errorf("checksum not found for %s", binaryName)
}