package config

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// Validate validates the configuration and returns any validation errors
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Validate general settings
	if c.DefaultProvider != "" {
		if !isValidProvider(c.DefaultProvider) {
			errors = append(errors, ValidationError{
				Field:   "default_provider",
				Message: fmt.Sprintf("invalid provider '%s', must be one of: openai, azure_openai, anthropic, gemini, deepseek", c.DefaultProvider),
			})
		}
	}


	// Validate provider configurations
	if c.OpenAI != nil {
		if err := validateProviderConfig("openai", c.OpenAI); err != nil {
			errors = append(errors, err...)
		}
	}

	if c.OpenAICompatible != nil {
		if err := validateProviderConfig("openai_compatible", c.OpenAICompatible); err != nil {
			errors = append(errors, err...)
		}
	}

	if c.AzureOpenAI != nil {
		if err := validateAzureOpenAIConfig(c.AzureOpenAI); err != nil {
			errors = append(errors, err...)
		}
	}

	if c.Anthropic != nil {
		if err := validateProviderConfig("anthropic", c.Anthropic); err != nil {
			errors = append(errors, err...)
		}
	}

	if c.Gemini != nil {
		if err := validateProviderConfig("gemini", c.Gemini); err != nil {
			errors = append(errors, err...)
		}
	}

	if c.DeepSeek != nil {
		if err := validateProviderConfig("deepseek", c.DeepSeek); err != nil {
			errors = append(errors, err...)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateProviderAvailable validates that the specified provider is configured and has an API key
func (c *Config) ValidateProviderAvailable(provider string) error {
	switch provider {
	case "openai":
		if c.OpenAI == nil {
			return fmt.Errorf("OpenAI provider not configured")
		}
		if c.OpenAI.APIKey == "" {
			return fmt.Errorf("OpenAI API key not configured")
		}
	case "openai_compatible":
		if c.OpenAICompatible == nil {
			return fmt.Errorf("OpenAI Compatible provider not configured")
		}
		if c.OpenAICompatible.APIKey == "" {
			return fmt.Errorf("OpenAI Compatible API key not configured")
		}
	case "azure_openai":
		if c.AzureOpenAI == nil {
			return fmt.Errorf("Azure OpenAI provider not configured")
		}
		if c.AzureOpenAI.APIKey == "" {
			return fmt.Errorf("Azure OpenAI API key not configured")
		}
		if c.AzureOpenAI.DeploymentName == "" {
			return fmt.Errorf("Azure OpenAI deployment name not configured")
		}
	case "anthropic":
		if c.Anthropic == nil {
			return fmt.Errorf("Anthropic provider not configured")
		}
		if c.Anthropic.APIKey == "" {
			return fmt.Errorf("Anthropic API key not configured")
		}
	case "gemini":
		if c.Gemini == nil {
			return fmt.Errorf("Gemini provider not configured")
		}
		if c.Gemini.APIKey == "" {
			return fmt.Errorf("Gemini API key not configured")
		}
	case "deepseek":
		if c.DeepSeek == nil {
			return fmt.Errorf("DeepSeek provider not configured")
		}
		if c.DeepSeek.APIKey == "" {
			return fmt.Errorf("DeepSeek API key not configured")
		}
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	return nil
}

// validateProviderConfig validates a basic provider configuration
func validateProviderConfig(providerName string, config *ProviderConfig) ValidationErrors {
	var errors ValidationErrors
	prefix := providerName

	// Validate base URL if provided
	if config.BaseURL != "" {
		if err := validateURL(config.BaseURL); err != nil {
			errors = append(errors, ValidationError{
				Field:   prefix + ".base_url",
				Message: err.Error(),
			})
		}
	}

	// Validate model name if provided
	if config.Model != "" {
		if err := validateModelName(providerName, config.Model); err != nil {
			errors = append(errors, ValidationError{
				Field:   prefix + ".model",
				Message: err.Error(),
			})
		}
	}

	return errors
}

// validateAzureOpenAIConfig validates Azure OpenAI specific configuration
func validateAzureOpenAIConfig(config *AzureOpenAIConfig) ValidationErrors {
	var errors ValidationErrors

	// Validate basic provider config
	errors = append(errors, validateProviderConfig("azure_openai", &config.ProviderConfig)...)

	// Validate Azure-specific fields
	if config.ResourceName != "" && config.BaseURL != "" {
		errors = append(errors, ValidationError{
			Field:   "azure_openai.resource_name",
			Message: "cannot specify both resource_name and base_url, use one or the other",
		})
	}

	if config.DeploymentName == "" && config.APIKey != "" {
		errors = append(errors, ValidationError{
			Field:   "azure_openai.deployment_name",
			Message: "deployment_name is required when using Azure OpenAI",
		})
	}

	if config.APIVersion != "" {
		if !isValidAzureAPIVersion(config.APIVersion) {
			errors = append(errors, ValidationError{
				Field:   "azure_openai.api_version",
				Message: "invalid API version format, should be in format YYYY-MM-DD",
			})
		}
	}

	return errors
}

// validateURL validates that a string is a valid URL
func validateURL(urlString string) error {
	if urlString == "" {
		return nil
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

// validateModelName validates provider-specific model names
func validateModelName(provider, model string) error {
	if model == "" {
		return nil
	}

	switch provider {
	case "openai":
		// OpenAI model validation
		validModels := []string{
			"gpt-4o", "gpt-4o-mini", "gpt-4", "gpt-4-turbo", "gpt-3.5-turbo",
			"gpt-4-32k", "gpt-4-0613", "gpt-4-32k-0613", "gpt-3.5-turbo-16k",
		}
		if !contains(validModels, model) && !strings.HasPrefix(model, "gpt-") {
			return fmt.Errorf("model '%s' may not be valid for OpenAI (expected format: gpt-*)", model)
		}
	case "anthropic":
		// Anthropic model validation
		if !strings.HasPrefix(model, "claude-") {
			return fmt.Errorf("model '%s' may not be valid for Anthropic (expected format: claude-*)", model)
		}
	case "gemini":
		// Gemini model validation
		if !strings.HasPrefix(model, "gemini-") && !strings.HasPrefix(model, "models/gemini-") {
			return fmt.Errorf("model '%s' may not be valid for Gemini (expected format: gemini-* or models/gemini-*)", model)
		}
	case "deepseek":
		// DeepSeek model validation
		if !strings.HasPrefix(model, "deepseek-") {
			return fmt.Errorf("model '%s' may not be valid for DeepSeek (expected format: deepseek-*)", model)
		}
	}

	return nil
}

// isValidProvider checks if the provider name is supported
func isValidProvider(provider string) bool {
	validProviders := []string{"openai", "openai_compatible", "azure_openai", "anthropic", "gemini", "deepseek"}
	return contains(validProviders, provider)
}

// isValidAzureAPIVersion validates Azure OpenAI API version format
func isValidAzureAPIVersion(version string) bool {
	// Basic format validation: YYYY-MM-DD
	if len(version) != 10 {
		return false
	}
	
	// Check format with simple pattern matching
	parts := strings.Split(version, "-")
	if len(parts) != 3 {
		return false
	}
	
	// Check year (4 digits)
	if len(parts[0]) != 4 {
		return false
	}
	
	// Check month (2 digits)
	if len(parts[1]) != 2 {
		return false
	}
	
	// Check day (2 digits)
	if len(parts[2]) != 2 {
		return false
	}
	
	return true
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}