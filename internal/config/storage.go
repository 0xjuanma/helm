package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDir  = ".helm"
	configFile = "settings.json"
)

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func ensureConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, configDir)
	return os.MkdirAll(dir, 0755)
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
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			cfg.Normalize()
			return cfg, nil
		}
		return nil, err
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

		// Migrate to Design workflow if it doesn't have sound
		if cfg.Design != nil && cfg.Design.Sound == nil {
			soundCopy := *legacy.Sound
			cfg.Design.Sound = &soundCopy
		}

		// Migrate to Custom workflow if it doesn't have sound
		if cfg.Custom != nil && cfg.Custom.Sound == nil {
			soundCopy := *legacy.Sound
			cfg.Custom.Sound = &soundCopy
		}
	}

	cfg.Normalize()

	// Save migrated config back to disk to clean up old format
	_ = Save(cfg)

	return cfg, nil
}

func Save(cfg *Config) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

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
