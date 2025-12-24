package state

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the overall configuration file
type Config struct {
	Data struct {
		Directory string
	}
	Convert struct {
		JSON bool
	}
}

// LoadConfig will load the config from disk
func LoadConfig() (Config, error) {
	cfg, err := detectDir("XDG_CONFIG_HOME", "CONFIG", ".config")
	if err != nil {
		return Config{}, err
	}
	cfg = filepath.Join(cfg, "settings.toml")
	config := Config{}
	if _, err := os.Stat(cfg); !errors.Is(err, os.ErrNotExist) {
		meta, err := toml.DecodeFile(cfg, &config)
		if err != nil {
			return config, err
		}

		undecoded := meta.Undecoded()
		if len(undecoded) > 0 {
			return config, fmt.Errorf("unknown config TOML fields: %v", undecoded)
		}
	}
	if config.Data.Directory == "" {
		d, err := detectDir("XDG_CACHE_HOME", "DATA_DIR", ".cache")
		if err != nil {
			return config, err
		}
		fmt.Println(d)
		config.Data.Directory = d
	}
	const isHome = "~"
	if strings.HasPrefix(config.Data.Directory, isHome) {
		home, err := os.UserHomeDir()
		if err != nil {
			return config, err
		}
		config.Data.Directory = strings.Replace(config.Data.Directory, isHome, home, 1)
	}
	return config, nil
}
