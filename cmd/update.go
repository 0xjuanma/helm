package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0xjuanma/helm/internal/version"
)

// runUpdate executes the appropriate update method based on installation detection.
// It first checks whether the running binary is already on the latest release (or
// is a dev build) and short-circuits with an informative message in that case.
func runUpdate() error {
	latest, fetchErr := version.CheckLatest()
	proceed, message := decideUpdate(Version, latest, fetchErr)
	if message != "" {
		fmt.Println(message)
	}
	if !proceed {
		return nil
	}

	installMethod := detectInstallationMethod()

	switch installMethod {
	case "homebrew":
		fmt.Println("Updating via Homebrew...")
		if err := runBrewUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "Homebrew update failed: %v\n", err)
			fmt.Println("Falling back to install script...")
			return runScriptUpdate()
		}
	default: // "script"
		fmt.Println("Updating via install script...")
		return runScriptUpdate()
	}
	return nil
}

// decideUpdate returns whether runUpdate should shell out to the installer, plus
// a user-facing message to print first. It is pure (no I/O) so the policy is unit
// testable without invoking brew or the install script.
//
// Policy:
//   - Dev builds short-circuit (cannot meaningfully self-update).
//   - On fetch failure we proceed: the user explicitly asked to update, so a
//     network flake shouldn't block them — fall back to today's behavior.
//   - If the current version is not older than latest (equal, or local ahead of
//     the published release), short-circuit with an informative message.
//   - Otherwise proceed.
func decideUpdate(current, latest string, fetchErr error) (proceed bool, message string) {
	if current == "dev" {
		return false, "Running a dev build; skipping update."
	}
	if fetchErr != nil || latest == "" {
		return true, ""
	}
	if version.IsOlder(current, latest) {
		return true, ""
	}
	return false, fmt.Sprintf("Already on the latest version (%s).", current)
}

// runBrewUpdate attempts to update helm via Homebrew.
func runBrewUpdate() error {
	cmd := exec.Command("brew", "upgrade", "0xjuanma/tap/helm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		// brew upgrade can exit non-zero for two recoverable reasons:
		//   1. The brew link step failed because a direct binary (e.g. from a
		//      prior script install) already exists at /usr/local/bin/helm,
		//      preventing Homebrew from creating its symlink.
		//   2. An unrelated brew cleanup error fires after a successful
		//      upgrade+link.
		// In both cases the formula was built successfully; attempt a forced
		// re-link before giving up and falling back to the script.
		fmt.Println("Attempting brew link recovery...")
		linkCmd := exec.Command("brew", "link", "--overwrite", "0xjuanma/tap/helm")
		linkCmd.Stdout = os.Stdout
		linkCmd.Stderr = os.Stderr
		if linkErr := linkCmd.Run(); linkErr == nil {
			return nil
		}
		return err
	}
	return nil
}

// runScriptUpdate updates helm via the install script.
func runScriptUpdate() error {
	cmd := exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/0xjuanma/helm/main/scripts/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// detectInstallationMethod returns "homebrew" or "script" based on how helm was installed.
func detectInstallationMethod() string {
	// 1. Fast path: check if binary is in Homebrew Cellar
	if isBinaryInCellar() {
		return "homebrew"
	}

	// 2. Fallback: ask brew directly if package is installed
	if isListedInBrew() {
		return "homebrew"
	}

	// 3. Default to script installation
	return "script"
}

// isBinaryInCellar checks if the helm binary is located in Homebrew's Cellar directory.
func isBinaryInCellar() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	realPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return false
	}

	return strings.Contains(realPath, "/Cellar/helm/")
}

// isListedInBrew checks if helm appears in brew's installed package list.
func isListedInBrew() bool {
	if _, err := exec.LookPath("brew"); err != nil {
		return false
	}

	cmd := exec.Command("brew", "list", "helm")
	return cmd.Run() == nil
}
