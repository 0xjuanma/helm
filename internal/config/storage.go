package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/0xjuanma/cli-toolkit/dirs"
)

const (
	appName    = "helm"
	configFile = "settings.json"
)

func configPath() (string, error) {
	dir, err := dirs.New(appName).ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

// readLegacyConfig reads settings from helm's pre-XDG location (~/.helm),
// used before the config store moved to dirs.ConfigDir(). Kept as a
// one-time fallback so existing users don't silently lose their config on
// Linux, where the new path differs from the old one.
func readLegacyConfig() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(home, ".helm", configFile))
}

// legacyConfig is used to unmarshal old config files that may have a global Sound field
type legacyConfig struct {
	Design             *WorkflowConfig `json:"design,omitempty"`
	Custom             *WorkflowConfig `json:"custom,omitempty"`
	TransitionDelaySec int             `json:"transition_delay_sec"`
	Sound              *SoundConfig    `json:"sound,omitempty"` // Legacy field, will be migrated
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		legacyData, legacyErr := readLegacyConfig()
		if legacyErr != nil {
			cfg := DefaultConfig()
			cfg.Normalize()
			return cfg, nil
		}
		data = legacyData
	}

	// First, try to unmarshal as legacy config to handle old format
	var legacy legacyConfig
	if err := json.Unmarshal(data, &legacy); err != nil {
		cfg := DefaultConfig()
		cfg.Normalize()
		return cfg, nil
	}

	// Migrate from legacy config to new format
	cfg := &Config{
		Design:             legacy.Design,
		Custom:             legacy.Custom,
		TransitionDelaySec: legacy.TransitionDelaySec,
	}

	// Migrate global Sound to workflows if they don't have sound config
	if legacy.Sound != nil {
		legacy.Sound.Normalize()

		migrateSound := func(wc *WorkflowConfig) {
			if wc != nil && wc.Sound == nil {
				soundCopy := *legacy.Sound
				wc.Sound = &soundCopy
			}
		}

		migrateSound(cfg.Design)
		migrateSound(cfg.Custom)
	}

	cfg.Normalize()

	// Save back to disk: cleans up old-format JSON, and migrates a
	// legacy-path config onto the new XDG-aware location.
	_ = Save(cfg)

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
