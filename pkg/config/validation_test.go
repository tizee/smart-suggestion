package config

import (
	"strings"
	"testing"
)

func TestHasVersionPrefix(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "URL ending with /v4/",
			url:      "https://open.bigmodel.cn/api/paas/v4/",
			expected: true,
		},
		{
			name:     "URL ending with /v4",
			url:      "https://open.bigmodel.cn/api/paas/v4",
			expected: true,
		},
		{
			name:     "URL ending with /v1",
			url:      "https://api.openai.com/v1",
			expected: true,
		},
		{
			name:     "URL ending with /v2",
			url:      "https://api.example.com/v2",
			expected: true,
		},
		{
			name:     "URL ending with /v10",
			url:      "https://api.example.com/v10",
			expected: true,
		},
		{
			name:     "URL with full chat/completions path",
			url:      "https://open.bigmodel.cn/api/paas/v4/chat/completions",
			expected: false,
		},
		{
			name:     "URL without version prefix",
			url:      "https://api.example.com/chat/completions",
			expected: false,
		},
		{
			name:     "URL with 'valid' as last segment",
			url:      "https://api.example.com/valid",
			expected: false,
		},
		{
			name:     "URL with 'version' as last segment",
			url:      "https://api.example.com/version",
			expected: false,
		},
		{
			name:     "URL with /v- (v followed by non-digit)",
			url:      "https://api.example.com/v-a",
			expected: false,
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasVersionPrefix(tt.url)
			if result != tt.expected {
				t.Errorf("hasVersionPrefix(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestValidateProviderConfig_OpenAICompatible(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		expectError bool
	}{
		{
			name:        "valid base URL without version prefix",
			baseURL:     "https://api.example.com",
			expectError: false,
		},
		{
			name:        "valid base URL with full path",
			baseURL:     "https://open.bigmodel.cn/api/paas/v4/chat/completions",
			expectError: false,
		},
		{
			name:        "invalid base URL ending with /v4",
			baseURL:     "https://open.bigmodel.cn/api/paas/v4/",
			expectError: true,
		},
		{
			name:        "invalid base URL ending with /v1",
			baseURL:     "https://api.example.com/v1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ProviderConfig{
				BaseURL: tt.baseURL,
			}
			errors := validateProviderConfig("openai_compatible", config)

			if tt.expectError {
				if len(errors) == 0 {
					t.Errorf("expected validation error for base_url=%s, but got none", tt.baseURL)
				}
				// Check that the error message contains helpful information
				for _, err := range errors {
					if err.Field == "openai_compatible.base_url" {
						if !strings.Contains(err.Message, "version prefix") {
							t.Errorf("error message should mention 'version prefix', got: %s", err.Message)
						}
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("expected no validation error for base_url=%s, but got: %v", tt.baseURL, errors)
				}
			}
		})
	}
}
