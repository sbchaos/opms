package config

// Package config is a set of types for interacting with the gh configuration files.
// Note: This package is intended for use only in gh, any other use cases are subject
// to breakage and non-backwards compatible updates.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	folderName = "opms"

	// custom directory
	configDir = "OPMS_CONFIG_DIR"

	// directory for linux/mac
	defaultConfig = ".config"
	defaultCache  = ".cache"

	// config directory in linux
	xdgConfigHome = "XDG_CONFIG_HOME"
	xdgCacheHome  = "XDG_CACHE_HOME"

	// appData is the config location for windows
	appData      = "AppData"
	localAppData = "LocalAppData"
)

var (
	cfg     *Config
	once    sync.Once
	loadErr = errors.New("unable to load config")
)

type Config struct {
	Version           string    `json:"version"`
	CurrentProfile    string    `json:"current_profile"`
	AvailableProfiles []Profile `json:"available_profiles"`
	//Aliases           *AliasConfig `json:"aliases"`
	mu sync.RWMutex
}

var Read = func(fallback *Config) (*Config, error) {
	once.Do(func() {
		cfg, loadErr = load(generalConfigFile(), fallback)
	})
	return cfg, loadErr
}

const defaultConfigStr = `
{
	"version": "1",
	"Editor": "vim",
	"aliases": {},

	"current_profile": "",
	"available_profiles": [],

	"auth": null,
}`

// ReadFromString takes a json string and returns a Config.
func ReadFromString(str string) *Config {
	var cfg1 *Config
	err := json.Unmarshal([]byte(str), &cfg1)
	if err != nil {
		return &Config{}
	}
	return cfg1
}

func DefaultConfig() *Config {
	return ReadFromString(defaultConfigStr)
}

func Write(c *Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = writeFile(generalConfigFile(), data)
	if err != nil {
		return err
	}

	return nil
}

func load(generalFilePath string, fallback *Config) (*Config, error) {
	generalConf, err := confFromFile(generalFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if generalConf != nil {
		return generalConf, nil
	}

	if fallback != nil {
		return fallback, nil
	}
	return &Config{}, nil
}

func generalConfigFile() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func confFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg1 *Config
	err = json.Unmarshal(data, &cfg1)
	if err != nil {
		return nil, err
	}
	return cfg1, nil
}

// ConfigDir path precedence: OPMS_CONFIG_DIR, XDG_CONFIG_HOME, AppData (windows only), HOME.
func ConfigDir() string {
	var path string
	if a := os.Getenv(configDir); a != "" {
		path = a
	} else if b := os.Getenv(xdgConfigHome); b != "" {
		path = filepath.Join(b, folderName)
	} else if c := os.Getenv(appData); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, strings.ToUpper(folderName))
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, defaultConfig, folderName)
	}
	return path
}

// CacheDir path precedence: XDG_CACHE_HOME, LocalAppData (windows only), HOME, temp.
func CacheDir() string {
	if a := os.Getenv(xdgCacheHome); a != "" {
		return filepath.Join(a, folderName)
	} else if b := os.Getenv(localAppData); runtime.GOOS == "windows" && b != "" {
		return filepath.Join(b, strings.ToUpper(folderName))
	} else if c, err := os.UserHomeDir(); err == nil {
		return filepath.Join(c, defaultCache, folderName)
	} else {
		// Note that this has a minor security issue because /tmp is world-writeable.
		return filepath.Join(os.TempDir(), "opms-cache")
	}
}

func writeFile(filename string, data []byte) (writeErr error) {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0600)
}

func (c *Config) GetCurrentProfile() *Profile {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.AvailableProfiles) < 1 {
		return &Profile{}
	}

	for i := 0; i < len(c.AvailableProfiles); i++ {
		if c.CurrentProfile == c.AvailableProfiles[i].Name {
			return &c.AvailableProfiles[i]
		}
	}

	return &c.AvailableProfiles[0]
}
