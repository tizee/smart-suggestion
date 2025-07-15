package privacy

import (
	"strings"
	"testing"
)

func TestDefaultFilterConfig(t *testing.T) {
	config := DefaultFilterConfig()
	
	if !config.Enabled {
		t.Error("Expected default config to be enabled")
	}
	
	if config.Level != FilterLevelBasic {
		t.Errorf("Expected default level to be Basic, got %v", config.Level)
	}
	
	if config.ReplacementText != "[REDACTED]" {
		t.Errorf("Expected default replacement text to be '[REDACTED]', got %s", config.ReplacementText)
	}
}

func TestNewFilter(t *testing.T) {
	config := &FilterConfig{
		Level:           FilterLevelBasic,
		Enabled:         true,
		ReplacementText: "***",
	}
	
	filter := NewFilter(config)
	
	if filter == nil {
		t.Error("Expected filter to be created")
	}
	
	if filter.config != config {
		t.Error("Expected filter config to match input config")
	}
}

func TestNewFilterWithNilConfig(t *testing.T) {
	filter := NewFilter(nil)
	
	if filter == nil {
		t.Error("Expected filter to be created with default config")
	}
	
	if !filter.config.Enabled {
		t.Error("Expected filter to use default enabled config")
	}
}

func TestFilterText_Disabled(t *testing.T) {
	config := &FilterConfig{
		Enabled: false,
		Level:   FilterLevelBasic,
	}
	
	filter := NewFilter(config)
	input := "export OPENAI_API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12"
	result := filter.FilterText(input)
	
	if result != input {
		t.Error("Expected no filtering when disabled")
	}
}

func TestFilterText_OpenAIAPIKey(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	testCases := []struct {
		name     string
		input    string
		expected bool // whether the input should be filtered
	}{
		{
			name:     "OpenAI API Key in export",
			input:    "export OPENAI_API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12",
			expected: true,
		},
		{
			name:     "OpenAI Project Key",
			input:    "export PROJECT_KEY=pk-1234567890abcdef1234567890abcdef1234567890abcdef12",
			expected: true,
		},
		{
			name:     "Regular text",
			input:    "ls -la /home/user",
			expected: false,
		},
		{
			name:     "API key in curl command",
			input:    "curl -H 'Authorization: Bearer sk-1234567890abcdef1234567890abcdef1234567890abcdef12'",
			expected: true,
		},
		{
			name:     "Anthropic API Key",
			input:    "export ANTHROPIC_API_KEY=sk-ant-1234567890abcdef1234567890abcdef",
			expected: true,
		},
		{
			name:     "Google API Key",
			input:    "GOOGLE_API_KEY=AIzaSyDxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			expected: true,
		},
		{
			name:     "AWS Access Key",
			input:    "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			expected: true,
		},
		{
			name:     "GitHub Token",
			input:    "set GITHUB_TOKEN=ghp_1234567890abcdef1234567890abcdef123456",
			expected: true,
		},
		{
			name:     "Stripe Secret Key",
			input:    "export STRIPE_SECRET_KEY=sk_test_1234567890abcdef1234567890abcdef",
			expected: true,
		},
		{
			name:     "Database URL",
			input:    "DATABASE_URL=postgres://user:pass@localhost:5432/dbname",
			expected: true,
		},
		{
			name:     "JWT Secret",
			input:    "export JWT_SECRET=super-secret-jwt-key-12345",
			expected: true,
		},
		{
			name:     "Custom Key with KEY suffix",
			input:    "export MY_CUSTOM_KEY=abc123def456ghi789",
			expected: true,
		},
		{
			name:     "Custom Token with TOKEN suffix",
			input:    "DEPLOYMENT_TOKEN=xyz789abc123def456",
			expected: true,
		},
		{
			name:     "Echo API Key command",
			input:    "echo $OPENAI_API_KEY",
			expected: true,
		},
		{
			name:     "Echo custom key command",
			input:    "echo $MY_SECRET_KEY",
			expected: true,
		},
		{
			name:     "Standalone secret output",
			input:    "sk-1234567890abcdef1234567890abcdef1234567890abcdef12",
			expected: true,
		},
		{
			name:     "GitHub token output",
			input:    "ghp_1234567890abcdef1234567890abcdef123456",
			expected: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.FilterText(tc.input)
			
			if tc.expected {
				if result == tc.input {
					t.Errorf("Expected input to be filtered, but it wasn't: %s", tc.input)
				}
				if !strings.Contains(result, "[REDACTED]") {
					t.Errorf("Expected result to contain [REDACTED], got: %s", result)
				}
			} else {
				if result != tc.input {
					t.Errorf("Expected input to remain unchanged, got: %s", result)
				}
			}
		})
	}
}

func TestFilterText_JWTToken(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	input := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	result := filter.FilterText(input)
	
	if result == input {
		t.Error("Expected JWT token to be filtered")
	}
	
	if !strings.Contains(result, "[REDACTED]") {
		t.Errorf("Expected result to contain [REDACTED], got: %s", result)
	}
}

func TestFilterText_DatabaseURL(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	testCases := []string{
		"mysql://user:password@localhost:3306/database",
		"postgresql://admin:secret123@db.example.com/mydb",
		"mongodb://user:pass@mongo.example.com:27017/app",
		"redis://user:password@redis.example.com:6379",
	}
	
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			result := filter.FilterText(input)
			
			if result == input {
				t.Errorf("Expected database URL to be filtered: %s", input)
			}
			
			if !strings.Contains(result, "[REDACTED]") {
				t.Errorf("Expected result to contain [REDACTED], got: %s", result)
			}
		})
	}
}

func TestFilterText_ModerateLevel(t *testing.T) {
	config := &FilterConfig{
		Level:           FilterLevelModerate,
		Enabled:         true,
		ReplacementText: "[HIDDEN]",
	}
	filter := NewFilter(config)
	
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "AWS Access Key",
			input:    "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			expected: true,
		},
		{
			name:     "GitHub Token",
			input:    "git remote set-url origin https://token:ghp_1234567890abcdef1234567890abcdef123456@github.com/user/repo.git",
			expected: true,
		},
		{
			name:     "Email in auth context",
			input:    "curl -u user@example.com:password123 https://api.example.com",
			expected: true,
		},
		{
			name:     "SSH Private Key",
			input:    "-----BEGIN RSA PRIVATE KEY-----",
			expected: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.FilterText(tc.input)
			
			if tc.expected {
				if result == tc.input {
					t.Errorf("Expected input to be filtered: %s", tc.input)
				}
				if !strings.Contains(result, "[HIDDEN]") {
					t.Errorf("Expected result to contain [HIDDEN], got: %s", result)
				}
			} else {
				if result != tc.input {
					t.Errorf("Expected input to remain unchanged, got: %s", result)
				}
			}
		})
	}
}

func TestFilterText_StrictLevel(t *testing.T) {
	config := &FilterConfig{
		Level:           FilterLevelStrict,
		Enabled:         true,
		ReplacementText: "***",
	}
	filter := NewFilter(config)
	
	// Test that strict level filters more aggressively
	input := "Here is a potential secret: abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
	result := filter.FilterText(input)
	
	if result == input {
		t.Error("Expected strict filtering to filter potential secrets")
	}
}

func TestFilterLines(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	lines := []string{
		"cd /home/user",
		"export API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12",
		"ls -la",
		"curl -H 'Authorization: Bearer token123' https://api.example.com",
	}
	
	result := filter.FilterLines(lines)
	
	if len(result) != len(lines) {
		t.Error("Expected same number of lines in result")
	}
	
	// First and third lines should be unchanged
	if result[0] != lines[0] || result[2] != lines[2] {
		t.Error("Expected non-sensitive lines to remain unchanged")
	}
	
	// Second and fourth lines should be filtered
	if result[1] == lines[1] {
		t.Error("Expected second line to be filtered")
	}
	if result[3] == lines[3] {
		t.Error("Expected fourth line to be filtered")
	}
}

func TestFilterMultilineText(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	input := `#!/bin/bash
cd /home/user
export OPENAI_API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12
curl -H "Authorization: Bearer $OPENAI_API_KEY" https://api.openai.com/v1/models
echo "Done"`
	
	result := filter.FilterMultilineText(input)
	
	if result == input {
		t.Error("Expected multiline text to be filtered")
	}
	
	if !strings.Contains(result, "[REDACTED]") {
		t.Error("Expected result to contain [REDACTED]")
	}
	
	// Check that non-sensitive lines are preserved
	if !strings.Contains(result, "#!/bin/bash") {
		t.Error("Expected shebang line to be preserved")
	}
	if !strings.Contains(result, "cd /home/user") {
		t.Error("Expected cd command to be preserved")
	}
}

func TestFilterText_EchoCommandAndOutput(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Echo command with API key",
			input:    "echo $OPENAI_API_KEY",
			expected: true,
		},
		{
			name:     "Echo command with custom key",
			input:    "echo $MY_SECRET_KEY",
			expected: true,
		},
		{
			name:     "Echo command output - OpenAI key",
			input:    "sk-1234567890abcdef1234567890abcdef1234567890abcdef12",
			expected: true,
		},
		{
			name:     "Echo command output - GitHub token",
			input:    "ghp_1234567890abcdef1234567890abcdef123456",
			expected: true,
		},
		{
			name:     "Echo command output - AWS key",
			input:    "AKIAIOSFODNN7EXAMPLE",
			expected: true,
		},
		{
			name:     "Echo normal text",
			input:    "echo 'Hello World'",
			expected: false,
		},
		{
			name:     "Normal command output",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "Terminal session with echo",
			input:    `$ echo $OPENAI_API_KEY
sk-1234567890abcdef1234567890abcdef1234567890abcdef12
$ ls -la`,
			expected: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.FilterMultilineText(tc.input)
			
			if tc.expected {
				if result == tc.input {
					t.Errorf("Expected input to be filtered: %s", tc.input)
				}
				if !strings.Contains(result, "[REDACTED]") {
					t.Errorf("Expected result to contain [REDACTED], got: %s", result)
				}
			} else {
				if result != tc.input {
					t.Errorf("Expected input to remain unchanged: %s -> %s", tc.input, result)
				}
			}
		})
	}
}

func TestDetectSensitivePatterns(t *testing.T) {
	filter := NewFilter(DefaultFilterConfig())
	
	input := "export OPENAI_API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12"
	detected := filter.DetectSensitivePatterns(input)
	
	if len(detected) == 0 {
		t.Error("Expected to detect sensitive patterns")
	}
	
	// Should detect both the export pattern and the OpenAI API key pattern
	expectedPatterns := []string{"Export API Key", "OpenAI API Key"}
	for _, expected := range expectedPatterns {
		found := false
		for _, detected := range detected {
			if strings.Contains(detected, "API Key") || strings.Contains(detected, "Export") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to detect pattern related to: %s, detected: %v", expected, detected)
		}
	}
}

func TestDetectSensitivePatterns_Disabled(t *testing.T) {
	config := &FilterConfig{
		Enabled: false,
		Level:   FilterLevelBasic,
	}
	filter := NewFilter(config)
	
	input := "export OPENAI_API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12"
	detected := filter.DetectSensitivePatterns(input)
	
	if len(detected) != 0 {
		t.Error("Expected no patterns to be detected when filter is disabled")
	}
}

func TestCustomPatterns(t *testing.T) {
	config := &FilterConfig{
		Level:           FilterLevelBasic,
		Enabled:         true,
		CustomPatterns:  []string{`my_secret_\w+`},
		ReplacementText: "[CUSTOM]",
	}
	
	filter := NewFilter(config)
	
	input := "export MY_VAR=my_secret_123456"
	result := filter.FilterText(input)
	
	if result == input {
		t.Error("Expected custom pattern to be filtered")
	}
	
	if !strings.Contains(result, "[CUSTOM]") {
		t.Errorf("Expected result to contain [CUSTOM], got: %s", result)
	}
}

func TestFilterLevels(t *testing.T) {
	testCases := []struct {
		level    FilterLevel
		input    string
		filtered bool
	}{
		{FilterLevelNone, "export API_KEY=sk-123", false},
		{FilterLevelBasic, "export API_KEY=sk-1234567890abcdef1234567890abcdef1234567890abcdef12", true},
		{FilterLevelModerate, "user@example.com", false}, // Email alone shouldn't be filtered
		{FilterLevelModerate, "export EMAIL=user@example.com", true}, // Email in export should be filtered
		{FilterLevelStrict, "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			config := &FilterConfig{
				Level:           tc.level,
				Enabled:         true,
				ReplacementText: "[FILTERED]",
			}
			filter := NewFilter(config)
			
			result := filter.FilterText(tc.input)
			
			if tc.filtered {
				if result == tc.input {
					t.Errorf("Expected input to be filtered at level %v: %s", tc.level, tc.input)
				}
			} else {
				if result != tc.input {
					t.Errorf("Expected input to remain unchanged at level %v: %s -> %s", tc.level, tc.input, result)
				}
			}
		})
	}
}