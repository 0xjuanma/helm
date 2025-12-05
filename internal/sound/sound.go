package sound

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/0xjuanma/helm/internal/config"
)

// Play emits a sound based on the provided configuration.
// If the system notification fails, it falls back to the terminal bell.
func Play(cfg config.SoundConfig) {
	cfg.Normalize()
	if !cfg.Enabled {
		return
	}

	if cfg.Mode == config.SoundModeMac && runtime.GOOS == "darwin" {
		tone := cfg.Tone
		if tone == "" {
			tone = config.DefaultMacTone
		}
		if err := playMacTone(tone); err == nil {
			return
		}
	}

	terminalBell()
}

func terminalBell() {
	_, _ = os.Stdout.Write([]byte{0x07})
	_, _ = os.Stderr.Write([]byte{0x07})
}

func playMacTone(name string) error {
	paths := []string{
		"/System/Library/Sounds/" + name + ".aiff",
		"/System/Library/Sounds/" + name + ".caf",
	}
	for _, path := range paths {
		cmd := exec.Command("afplay", path)
		if err := cmd.Start(); err != nil {
			continue
		}
		return cmd.Wait()
	}
	return os.ErrNotExist
}
