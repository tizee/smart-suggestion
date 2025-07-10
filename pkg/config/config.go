package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProviderConfig represents the configuration for a single AI provider
type ProviderConfig struct {
	APIKey     string `json:"api_key,omitempty"`
	BaseURL    string `json:"base_url,omitempty"`
	Model      string `json:"model,omitempty"`
	APIVersion string `json:"api_version,omitempty"`
}

// AzureOpenAIConfig represents specific configuration for Azure OpenAI
type AzureOpenAIConfig struct {
	ProviderConfig
	ResourceName   string `json:"resource_name,omitempty"`
	DeploymentName string `json:"deployment_name,omitempty"`
}

// Config represents the complete application configuration
type Config struct {
	// Provider configurations
	OpenAI           *ProviderConfig    `json:"openai,omitempty"`
	OpenAICompatible *ProviderConfig    `json:"openai_compatible,omitempty"`
	AzureOpenAI      *AzureOpenAIConfig `json:"azure_openai,omitempty"`
	Anthropic        *ProviderConfig    `json:"anthropic,omitempty"`
	Gemini           *ProviderConfig    `json:"gemini,omitempty"`
	DeepSeek         *ProviderConfig    `json:"deepseek,omitempty"`

	// General settings (only provider-related)
	DefaultProvider string `json:"default_provider,omitempty"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		DefaultProvider: "openai",
		OpenAI: &ProviderConfig{
			BaseURL: "https://api.openai.com",
			Model:   "gpt-4o-mini",
		},
		OpenAICompatible: &ProviderConfig{
			BaseURL: "http://localhost:11434",
			Model:   "llama3.2:latest",
		},
		AzureOpenAI: &AzureOpenAIConfig{
			ProviderConfig: ProviderConfig{
				APIVersion: "2024-10-21",
			},
		},
		Anthropic: &ProviderConfig{
			BaseURL: "https://api.anthropic.com",
			Model:   "claude-3-5-sonnet-20241022",
		},
		Gemini: &ProviderConfig{
			BaseURL: "https://generativelanguage.googleapis.com",
			Model:   "gemini-2.5-flash",
		},
		DeepSeek: &ProviderConfig{
			BaseURL: "https://api.deepseek.com",
			Model:   "deepseek-chat",
		},
	}
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "smart-suggestion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads configuration from the specified file path
// If the file doesn't exist, returns an error
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config file path is required")
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults for missing values
	defaultConfig := DefaultConfig()
	mergeConfigs(&config, defaultConfig)

	return &config, nil
}

// LoadConfigFromEnv loads configuration from the path specified in SMART_SUGGESTION_PROVIDER_FILE
// environment variable. If the environment variable is not set, returns an error.
func LoadConfigFromEnv() (*Config, error) {
	configPath := os.Getenv("SMART_SUGGESTION_PROVIDER_FILE")
	if configPath == "" {
		return nil, fmt.Errorf("SMART_SUGGESTION_PROVIDER_FILE environment variable is not set")
	}

	return LoadConfig(configPath)
}

// SaveConfig saves the configuration to the specified file path
func (c *Config) SaveConfig(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("config file path is required")
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetProviderConfig returns the configuration for the specified provider
func (c *Config) GetProviderConfig(provider string) (*ProviderConfig, error) {
	switch provider {
	case "openai":
		if c.OpenAI == nil {
			return nil, fmt.Errorf("OpenAI configuration not found")
		}
		return c.OpenAI, nil
	case "openai_compatible":
		if c.OpenAICompatible == nil {
			return nil, fmt.Errorf("OpenAI Compatible configuration not found")
		}
		return c.OpenAICompatible, nil
	case "anthropic":
		if c.Anthropic == nil {
			return nil, fmt.Errorf("Anthropic configuration not found")
		}
		return c.Anthropic, nil
	case "gemini":
		if c.Gemini == nil {
			return nil, fmt.Errorf("Gemini configuration not found")
		}
		return c.Gemini, nil
	case "deepseek":
		if c.DeepSeek == nil {
			return nil, fmt.Errorf("DeepSeek configuration not found")
		}
		return c.DeepSeek, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// GetAzureOpenAIConfig returns the Azure OpenAI configuration
func (c *Config) GetAzureOpenAIConfig() (*AzureOpenAIConfig, error) {
	if c.AzureOpenAI == nil {
		return nil, fmt.Errorf("Azure OpenAI configuration not found")
	}
	return c.AzureOpenAI, nil
}

// GetAPIKey gets the API key from config only (no environment variable fallback)
func (c *Config) GetAPIKey(provider string) (string, error) {
	var configKey string

	// Get API key from config
	switch provider {
	case "openai":
		if c.OpenAI != nil {
			configKey = c.OpenAI.APIKey
		}
	case "openai_compatible":
		if c.OpenAICompatible != nil {
			configKey = c.OpenAICompatible.APIKey
		}
	case "azure_openai":
		if c.AzureOpenAI != nil {
			configKey = c.AzureOpenAI.APIKey
		}
	case "anthropic":
		if c.Anthropic != nil {
			configKey = c.Anthropic.APIKey
		}
	case "gemini":
		if c.Gemini != nil {
			configKey = c.Gemini.APIKey
		}
	case "deepseek":
		if c.DeepSeek != nil {
			configKey = c.DeepSeek.APIKey
		}
	}

	// Return config key if available
	if configKey != "" {
		return configKey, nil
	}

	return "", fmt.Errorf("%s API key not found in config file", provider)
}

// mergeConfigs merges missing fields from defaultConfig into config
func mergeConfigs(config, defaultConfig *Config) {
	if config.DefaultProvider == "" {
		config.DefaultProvider = defaultConfig.DefaultProvider
	}

	// Merge provider configs
	if config.OpenAI == nil {
		config.OpenAI = defaultConfig.OpenAI
	} else {
		mergeProviderConfig(config.OpenAI, defaultConfig.OpenAI)
	}

	if config.OpenAICompatible == nil {
		config.OpenAICompatible = defaultConfig.OpenAICompatible
	} else {
		mergeProviderConfig(config.OpenAICompatible, defaultConfig.OpenAICompatible)
	}

	if config.AzureOpenAI == nil {
		config.AzureOpenAI = defaultConfig.AzureOpenAI
	} else {
		mergeProviderConfig(&config.AzureOpenAI.ProviderConfig, &defaultConfig.AzureOpenAI.ProviderConfig)
		if config.AzureOpenAI.APIVersion == "" {
			config.AzureOpenAI.APIVersion = defaultConfig.AzureOpenAI.APIVersion
		}
	}

	if config.Anthropic == nil {
		config.Anthropic = defaultConfig.Anthropic
	} else {
		mergeProviderConfig(config.Anthropic, defaultConfig.Anthropic)
	}

	if config.Gemini == nil {
		config.Gemini = defaultConfig.Gemini
	} else {
		mergeProviderConfig(config.Gemini, defaultConfig.Gemini)
	}

	if config.DeepSeek == nil {
		config.DeepSeek = defaultConfig.DeepSeek
	} else {
		mergeProviderConfig(config.DeepSeek, defaultConfig.DeepSeek)
	}
}

// mergeProviderConfig merges missing fields from defaultProvider into provider
func mergeProviderConfig(provider, defaultProvider *ProviderConfig) {
	if provider.BaseURL == "" {
		provider.BaseURL = defaultProvider.BaseURL
	}
	if provider.Model == "" {
		provider.Model = defaultProvider.Model
	}
	if provider.APIVersion == "" {
		provider.APIVersion = defaultProvider.APIVersion
	}
}