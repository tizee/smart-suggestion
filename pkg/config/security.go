package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetSecureFilePermissions sets restrictive permissions on a file (0600 - owner read/write only)
func SetSecureFilePermissions(filepath string) error {
	return os.Chmod(filepath, 0600)
}

// CreateSecureDirectory creates a directory with secure permissions (0700 - owner access only)
func CreateSecureDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0700)
}

// SecureConfigPath returns a secure default path for configuration files
func SecureConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Use .config directory with secure permissions
	configDir := filepath.Join(homeDir, ".config", "smart-suggestion")
	if err := CreateSecureDirectory(configDir); err != nil {
		return "", fmt.Errorf("failed to create secure config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	return configPath, nil
}