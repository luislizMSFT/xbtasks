package ado

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// AzCliTimeout is the maximum time to wait for an az CLI command to complete.
const AzCliTimeout = 15 * time.Second

// ResolveAzPath finds the az CLI binary, expanding PATH for macOS .app bundles.
func ResolveAzPath() string {
	if p, err := exec.LookPath("az"); err == nil {
		return p
	}
	if runtime.GOOS == "darwin" {
		for _, candidate := range []string{
			"/opt/homebrew/bin/az",
			"/usr/local/bin/az",
		} {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}
	return "az"
}

// RunAzCli executes an az CLI command with a default timeout and returns the output.
func RunAzCli(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AzCliTimeout)
	defer cancel()
	return RunAzCliCtx(ctx, args...)
}

// RunAzCliCtx executes an az CLI command with the given context and returns the output.
func RunAzCliCtx(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, ResolveAzPath(), args...)
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("az cli timed out after %s", AzCliTimeout)
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("az cli error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("az cli failed: %w", err)
	}
	return output, nil
}

// EnsurePATH adds common tool directories for macOS .app bundles.
func EnsurePATH() {
	if runtime.GOOS != "darwin" {
		return
	}
	extra := []string{"/opt/homebrew/bin", "/usr/local/bin"}
	cur := os.Getenv("PATH")
	var missing []string
	for _, d := range extra {
		if !strings.Contains(cur, d) {
			missing = append(missing, d)
		}
	}
	if len(missing) > 0 {
		os.Setenv("PATH", cur+":"+strings.Join(missing, ":"))
	}
}

func init() {
	EnsurePATH()
}
