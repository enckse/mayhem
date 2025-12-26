package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const databaseName = FileName + "db"

// Config is the overall configuration file
type Config struct {
	Data struct {
		Directory string
		Pretty    bool
	}
	Display struct {
		Finished struct {
			Since string
		}
	}
	Backups struct {
		Directory string
		Format    string
		Duration  string
	}
}

// Database will get the path to the database file
func (c Config) Database() string {
	return filepath.Join(c.Data.Directory, databaseName)
}

// LoadConfig will load the config from disk
func LoadConfig(file string) (Config, error) {
	cfg := file
	if cfg == "" {
		var err error
		cfg, err = detectDir("XDG_CONFIG_HOME", "CONFIG", ".config")
		if err != nil {
			return Config{}, err
		}
		cfg = filepath.Join(cfg, "settings.toml")
	}
	config := Config{}
	if PathExists(cfg) {
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
	if config.Backups.Directory != "" {
		config.Backups.Directory = filepath.Join(config.Data.Directory, config.Backups.Directory)
	}
	return config, nil
}
