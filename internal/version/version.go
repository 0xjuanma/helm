// Package version provides helpers for comparing and fetching helm release versions.
package version

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// IsOlder returns true if versionA is older than versionB.
// Supports semantic versions like "v1.2.3" or "1.2.3".
func IsOlder(versionA, versionB string) bool {
	// Remove 'v' prefix for comparison
	va := strings.TrimPrefix(versionA, "v")
	vb := strings.TrimPrefix(versionB, "v")

	// Split by dots
	partsA := strings.Split(va, ".")
	partsB := strings.Split(vb, ".")

	// Compare each part numerically
	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		numA, errA := strconv.Atoi(partsA[i])
		numB, errB := strconv.Atoi(partsB[i])

		if errA != nil || errB != nil {
			// If parsing fails, fall back to string comparison
			return va < vb
		}

		if numA < numB {
			return true
		}
		if numA > numB {
			return false
		}
	}

	// If all compared parts are equal, longer version is considered newer
	return len(partsA) < len(partsB)
}

// CheckLatest fetches the latest release version from GitHub.
// Uses GitHub's redirect URL which is simpler than the API.
// Returns the version tag (e.g., "v1.2.3").
func CheckLatest() (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://github.com/0xjuanma/helm/releases/latest")
	if err != nil {
		return "", fmt.Errorf("fetch latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// GitHub redirects to: https://github.com/0xjuanma/helm/releases/tag/v1.2.3
	// Extract version from the final URL
	finalURL := resp.Request.URL.String()

	if idx := strings.LastIndex(finalURL, "/releases/tag/"); idx != -1 {
		version := finalURL[idx+len("/releases/tag/"):]
		if version != "" {
			return version, nil
		}
	}

	// Fallback: try to read from response body (in case redirect doesn't work)
	body, err := io.ReadAll(resp.Body)
	if err == nil && len(body) > 0 {
		bodyStr := string(body)
		if idx := strings.Index(bodyStr, "/releases/tag/"); idx != -1 {
			start := idx + len("/releases/tag/")
			end := strings.Index(bodyStr[start:], "\"")
			if end != -1 {
				version := bodyStr[start : start+end]
				if version != "" {
					return version, nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not extract version from GitHub response")
}
