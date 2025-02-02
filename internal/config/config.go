package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    StoragePath string
}

func Load() (*Config, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        configDir = os.Getenv("HOME")
    }

    return &Config{
        StoragePath: filepath.Join(configDir, "directory-bookmarks.json"),
    }, nil
}
