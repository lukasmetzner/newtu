package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	RSSFeeds   []RSSFeed `json:"rss_feeds"`
	DBFilePath string    `json:"-"`
}

type RSSFeed struct {
	Source string `json:"source"`
	URL    string `json:"url"`
}

func (c *Config) Save() error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

func ReadConfig() (*Config, error) {
	configDirPath, err := configDir()
	if err != nil {
		return nil, err
	}

	cacheDirPath, err := cacheDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(configDirPath, "config.json"))
	if os.IsNotExist(err) {
		cfg := &Config{
			DBFilePath: filepath.Join(cacheDirPath, "data.db"),
		}
		if saveErr := cfg.Save(); saveErr != nil {
			return nil, fmt.Errorf("creating default config: %w", saveErr)
		}
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	cfg.DBFilePath = filepath.Join(cacheDirPath, "data.db")

	return &cfg, nil
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting config directory: %w", err)
	}

	dir := filepath.Join(base, "newtu")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	return dir, nil
}

func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("getting cache directory: %w", err)
	}

	dir := filepath.Join(base, "newtu")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating cache directory: %w", err)
	}

	return dir, nil
}
