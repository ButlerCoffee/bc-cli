package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DefaultAPIURL = "https://api.butler.coffee"
	ConfigDir     = ".butler-coffee"
	ConfigFile    = "config.json"
)

type Config struct {
	APIURL       string `json:"api_url"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GetAPIURL() string {
	// Check for BASE_HOSTNAME environment variable
	if hostname := os.Getenv("BASE_HOSTNAME"); hostname != "" {
		return hostname
	}
	return DefaultAPIURL
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	apiURL := GetAPIURL()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{APIURL: apiURL}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variable if set, otherwise use config or default
	if envURL := os.Getenv("BASE_HOSTNAME"); envURL != "" {
		cfg.APIURL = envURL
	} else if cfg.APIURL == "" {
		cfg.APIURL = apiURL
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != ""
}
