package privacy

import (
	"regexp"
	"strings"
)

// FilterLevel represents the sensitivity level of privacy filtering
type FilterLevel int

const (
	// FilterLevelNone disables all privacy filtering
	FilterLevelNone FilterLevel = iota
	// FilterLevelBasic enables basic filtering for common sensitive patterns
	FilterLevelBasic
	// FilterLevelModerate enables moderate filtering including emails and IPs
	FilterLevelModerate
	// FilterLevelStrict enables strict filtering including aggressive pattern matching
	FilterLevelStrict
)

// FilterConfig represents the configuration for privacy filtering
type FilterConfig struct {
	Level           FilterLevel `json:"level"`
	Enabled         bool        `json:"enabled"`
	CustomPatterns  []string    `json:"custom_patterns,omitempty"`
	ReplacementText string      `json:"replacement_text,omitempty"`
}

// DefaultFilterConfig returns a default privacy filter configuration
func DefaultFilterConfig() *FilterConfig {
	return &FilterConfig{
		Level:           FilterLevelBasic,
		Enabled:         true,
		CustomPatterns:  []string{},
		ReplacementText: "[REDACTED]",
	}
}

// SensitivePattern represents a pattern to detect sensitive information
type SensitivePattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
	Level       FilterLevel
}

// Filter represents the privacy filter with compiled patterns
type Filter struct {
	config   *FilterConfig
	patterns []SensitivePattern
}

// NewFilter creates a new privacy filter with the given configuration
func NewFilter(config *FilterConfig) *Filter {
	if config == nil {
		config = DefaultFilterConfig()
	}

	filter := &Filter{
		config:   config,
		patterns: []SensitivePattern{},
	}

	filter.compilePatterns()
	return filter
}

// compilePatterns compiles all the sensitive patterns based on the filter level
func (f *Filter) compilePatterns() {
	replacementText := f.config.ReplacementText
	if replacementText == "" {
		replacementText = "[REDACTED]"
	}

	// Basic level patterns - common API keys and tokens
	basicPatterns := []struct {
		name    string
		pattern string
	}{
		// OpenAI API keys
		{"OpenAI API Key", `sk-[a-zA-Z0-9]{48,}`},
		{"OpenAI Project Key", `pk-[a-zA-Z0-9]{48,}`},
		
		// Common API key patterns
		{"Generic API Key", `(?i)api[_-]?key['"=:\s]+['"]*([a-zA-Z0-9_\-]{8,})['"]*`},
		{"Bearer Token", `(?i)bearer\s+([a-zA-Z0-9_\-\.]{2,})`},
		{"Authorization Header", `(?i)authorization['"=:\s]+['"]*([a-zA-Z0-9_\-\.]{2,})['"]*`},
		
		// Environment variable exports containing secrets
		{"Export API Key", `(?i)export\s+[A-Z_]*(?:API|KEY|TOKEN|SECRET|PASSWORD)[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		{"Set Environment", `(?i)set\s+[A-Z_]*(?:API|KEY|TOKEN|SECRET|PASSWORD)[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		
		// Environment variable names containing KEY (broader pattern)
		{"Env Var with KEY", `(?i)(?:export\s+|set\s+)?[A-Z_]*KEY[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		{"Env Var with TOKEN", `(?i)(?:export\s+|set\s+)?[A-Z_]*TOKEN[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		{"Env Var with SECRET", `(?i)(?:export\s+|set\s+)?[A-Z_]*SECRET[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		{"Env Var with PASSWORD", `(?i)(?:export\s+|set\s+)?[A-Z_]*PASSWORD[A-Z_]*=['"]*([^'"'\s]{8,})['"]*`},
		
		// Echo command outputs that reveal secrets
		{"Echo API Key", `(?i)echo\s+\$[A-Z_]*(?:API|KEY|TOKEN|SECRET|PASSWORD)[A-Z_]*`},
		{"Echo Env Var", `(?i)echo\s+\$[A-Z_]*(?:KEY|TOKEN|SECRET|PASSWORD)[A-Z_]*`},
		
		// Command substitution outputs
		{"Command Substitution Secret", `(?i)\$\([^)]*(?:API|KEY|TOKEN|SECRET|PASSWORD)[^)]*\)`},
		
		// Standalone secret values that might be command outputs
		{"Standalone Secret Value", `(?m)^[a-zA-Z0-9_\-\.+/=]{20,}$`},
		
		// Lines that look like they contain revealed secrets (common patterns)
		{"Revealed Secret Line", `(?i)(?:^|\s)(?:sk-[a-zA-Z0-9]{48,}|pk-[a-zA-Z0-9]{48,}|ghp_[a-zA-Z0-9]{36}|ghs_[a-zA-Z0-9]{36}|AKIA[0-9A-Z]{16}|xox[baprs]-[0-9a-zA-Z\-]{10,72})(?:\s|$)`},
		
		// Common API key environment variable patterns
		{"OpenAI API Key Env", `(?i)(?:export\s+|set\s+)?OPENAI_API_KEY=['"]*([^'"'\s]{8,})['"]*`},
		{"Anthropic API Key Env", `(?i)(?:export\s+|set\s+)?ANTHROPIC_API_KEY=['"]*([^'"'\s]{8,})['"]*`},
		{"Google API Key Env", `(?i)(?:export\s+|set\s+)?(?:GOOGLE_API_KEY|GEMINI_API_KEY)=['"]*([^'"'\s]{8,})['"]*`},
		{"AWS Keys Env", `(?i)(?:export\s+|set\s+)?(?:AWS_ACCESS_KEY_ID|AWS_SECRET_ACCESS_KEY)=['"]*([^'"'\s]{8,})['"]*`},
		{"GitHub Token Env", `(?i)(?:export\s+|set\s+)?(?:GITHUB_TOKEN|GH_TOKEN)=['"]*([^'"'\s]{8,})['"]*`},
		{"Azure Keys Env", `(?i)(?:export\s+|set\s+)?(?:AZURE_CLIENT_SECRET|AZURE_TENANT_ID)=['"]*([^'"'\s]{8,})['"]*`},
		{"Slack Token Env", `(?i)(?:export\s+|set\s+)?(?:SLACK_TOKEN|SLACK_BOT_TOKEN)=['"]*([^'"'\s]{8,})['"]*`},
		{"DeepSeek API Key Env", `(?i)(?:export\s+|set\s+)?DEEPSEEK_API_KEY=['"]*([^'"'\s]{8,})['"]*`},
		{"Stripe Keys Env", `(?i)(?:export\s+|set\s+)?(?:STRIPE_SECRET_KEY|STRIPE_PUBLISHABLE_KEY)=['"]*([^'"'\s]{8,})['"]*`},
		{"Twilio Keys Env", `(?i)(?:export\s+|set\s+)?(?:TWILIO_AUTH_TOKEN|TWILIO_ACCOUNT_SID)=['"]*([^'"'\s]{8,})['"]*`},
		{"SendGrid API Key Env", `(?i)(?:export\s+|set\s+)?SENDGRID_API_KEY=['"]*([^'"'\s]{8,})['"]*`},
		{"Mailgun API Key Env", `(?i)(?:export\s+|set\s+)?MAILGUN_API_KEY=['"]*([^'"'\s]{8,})['"]*`},
		{"Redis URL Env", `(?i)(?:export\s+|set\s+)?REDIS_URL=['"]*([^'"'\s]{8,})['"]*`},
		{"MongoDB URI Env", `(?i)(?:export\s+|set\s+)?(?:MONGODB_URI|MONGO_URL)=['"]*([^'"'\s]{8,})['"]*`},
		{"Database URL Env", `(?i)(?:export\s+|set\s+)?(?:DATABASE_URL|DB_URL)=['"]*([^'"'\s]{8,})['"]*`},
		{"JWT Secret Env", `(?i)(?:export\s+|set\s+)?(?:JWT_SECRET|JWT_KEY)=['"]*([^'"'\s]{8,})['"]*`},
		{"Encryption Key Env", `(?i)(?:export\s+|set\s+)?(?:ENCRYPTION_KEY|SECRET_KEY|SESSION_SECRET)=['"]*([^'"'\s]{8,})['"]*`},
		{"Docker Registry Env", `(?i)(?:export\s+|set\s+)?(?:DOCKER_PASSWORD|REGISTRY_TOKEN)=['"]*([^'"'\s]{8,})['"]*`},
		{"CI/CD Token Env", `(?i)(?:export\s+|set\s+)?(?:CI_TOKEN|GITLAB_TOKEN|JENKINS_TOKEN)=['"]*([^'"'\s]{8,})['"]*`},
		{"Cloud Provider Keys", `(?i)(?:export\s+|set\s+)?(?:DIGITALOCEAN_TOKEN|VULTR_API_KEY|LINODE_TOKEN)=['"]*([^'"'\s]{8,})['"]*`},
		
		// JWT tokens
		{"JWT Token", `eyJ[a-zA-Z0-9_\-]*\.eyJ[a-zA-Z0-9_\-]*\.[a-zA-Z0-9_\-]*`},
		
		// Common secret patterns in command line
		{"Password Parameter", `(?i)--password[=\s]+['"]*([^'"'\s]{4,})['"]*`},
		{"Token Parameter", `(?i)--token[=\s]+['"]*([^'"'\s]{8,})['"]*`},
		{"Secret Parameter", `(?i)--secret[=\s]+['"]*([^'"'\s]{8,})['"]*`},
		
		// Database connection strings
		{"Database URL", `(?i)(mysql|postgresql|mongodb|redis)://[^@]+:[^@]+@[^\s]+`},
		
		// Generic secrets in curl/wget commands
		{"Curl Header Secret", `(?i)curl[^|]*-H['"]*[^'"]*(?:authorization|api[_-]?key|token)['"]*[=:]['"]*([^'"'\s]{8,})['"]*`},
		{"Wget Header Secret", `(?i)wget[^|]*--header[='"]*[^'"]*(?:authorization|api[_-]?key|token)['"]*[=:]['"]*([^'"'\s]{8,})['"]*`},
	}

	// Add basic patterns
	for _, p := range basicPatterns {
		if compiled, err := regexp.Compile(p.pattern); err == nil {
			f.patterns = append(f.patterns, SensitivePattern{
				Name:        p.name,
				Pattern:     compiled,
				Replacement: replacementText,
				Level:       FilterLevelBasic,
			})
		}
	}

	// Moderate level patterns - emails, IPs, more aggressive patterns
	if f.config.Level >= FilterLevelModerate {
		moderatePatterns := []struct {
			name    string
			pattern string
		}{
			// Email addresses in sensitive contexts
			{"Email in Auth", `(?i)(?:user|username|email|login)['"=:\s]+['"]*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})['"]*`},
			{"Email in curl -u", `(?i)curl\s+[^|]*-u\s+([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}):([^@\s]+)`},
			
			// IP addresses in sensitive contexts
			{"Private IP", `(?:192\.168\.|10\.|172\.(?:1[6-9]|2[0-9]|3[01])\.)\d{1,3}\.\d{1,3}(?::\d+)?`},
			
			// SSH private key patterns
			{"SSH Private Key", `-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----`},
			
			// AWS keys
			{"AWS Access Key", `AKIA[0-9A-Z]{16}`},
			{"AWS Secret Key", `(?i)aws[_-]?secret[_-]?access[_-]?key['"=:\s]+['"]*([a-zA-Z0-9/+]{40})['"]*`},
			
			// GitHub tokens
			{"GitHub Token", `ghp_[a-zA-Z0-9]{36}`},
			{"GitHub App Token", `ghs_[a-zA-Z0-9]{36}`},
			{"GitHub OAuth Token", `gho_[a-zA-Z0-9]{36}`},
			
			// Slack tokens
			{"Slack Token", `xox[baprs]-[0-9a-zA-Z-]{10,72}`},
			
			// More aggressive password detection
			{"Password in URL", `(?i)://[^:@]+:([^@\s]{4,})@`},
		}

		for _, p := range moderatePatterns {
			if compiled, err := regexp.Compile(p.pattern); err == nil {
				f.patterns = append(f.patterns, SensitivePattern{
					Name:        p.name,
					Pattern:     compiled,
					Replacement: replacementText,
					Level:       FilterLevelModerate,
				})
			}
		}
	}

	// Strict level patterns - very aggressive filtering
	if f.config.Level >= FilterLevelStrict {
		strictPatterns := []struct {
			name    string
			pattern string
		}{
			// Any long alphanumeric strings that could be secrets
			{"Potential Secret", `\b[a-zA-Z0-9]{32,}\b`},
			
			// Credit card numbers
			{"Credit Card", `\b(?:4\d{3}|5[1-5]\d{2}|6011|65\d{2})\s*\d{4}\s*\d{4}\s*\d{4}\b`},
			
			// Social Security Numbers (US format)
			{"SSN", `\b\d{3}-\d{2}-\d{4}\b`},
			
			// Phone numbers in sensitive contexts
			{"Phone Number", `(?i)(?:phone|tel|mobile)['"=:\s]+['"]*([+]?[\d\s\-\(\)]{10,})['"]*`},
		}

		for _, p := range strictPatterns {
			if compiled, err := regexp.Compile(p.pattern); err == nil {
				f.patterns = append(f.patterns, SensitivePattern{
					Name:        p.name,
					Pattern:     compiled,
					Replacement: replacementText,
					Level:       FilterLevelStrict,
				})
			}
		}
	}

	// Add custom patterns
	for _, customPattern := range f.config.CustomPatterns {
		if compiled, err := regexp.Compile(customPattern); err == nil {
			f.patterns = append(f.patterns, SensitivePattern{
				Name:        "Custom Pattern",
				Pattern:     compiled,
				Replacement: replacementText,
				Level:       FilterLevelBasic,
			})
		}
	}
}

// FilterText filters sensitive information from the given text
func (f *Filter) FilterText(text string) string {
	if !f.config.Enabled || f.config.Level == FilterLevelNone {
		return text
	}

	filtered := text

	// Apply each pattern
	for _, pattern := range f.patterns {
		if pattern.Level <= f.config.Level {
			filtered = pattern.Pattern.ReplaceAllString(filtered, pattern.Replacement)
		}
	}

	return filtered
}

// FilterLines filters sensitive information from multiple lines of text
func (f *Filter) FilterLines(lines []string) []string {
	if !f.config.Enabled || f.config.Level == FilterLevelNone {
		return lines
	}

	filtered := make([]string, len(lines))
	for i, line := range lines {
		filtered[i] = f.FilterText(line)
	}

	return filtered
}

// FilterMultilineText filters sensitive information from multiline text
func (f *Filter) FilterMultilineText(text string) string {
	if !f.config.Enabled || f.config.Level == FilterLevelNone {
		return text
	}

	lines := strings.Split(text, "\n")
	filteredLines := f.FilterLines(lines)
	return strings.Join(filteredLines, "\n")
}

// DetectSensitivePatterns returns information about detected sensitive patterns without filtering
func (f *Filter) DetectSensitivePatterns(text string) []string {
	if !f.config.Enabled || f.config.Level == FilterLevelNone {
		return []string{}
	}

	var detected []string

	for _, pattern := range f.patterns {
		if pattern.Level <= f.config.Level && pattern.Pattern.MatchString(text) {
			detected = append(detected, pattern.Name)
		}
	}

	return detected
}