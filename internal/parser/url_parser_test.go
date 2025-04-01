package parser

import (
	"reflect"
	"testing"
)

func TestURLParserParse(t *testing.T) {
	parser := NewURLParser()

	tests := []struct {
		name        string
		input       string
		expected    *URLParseResult
		expectError bool
	}{
		{
			name:  "Valid URL with HTTPS",
			input: "https://example.com/path?param=value",
			expected: &URLParseResult{
				OriginalURL: "https://example.com/path?param=value",
				Normalized:  "https://example.com/path?param=value",
				Domain:      "example.com",
				Path:        "/path",
				Params:      map[string]string{"param": "value"},
				IsValid:     true,
			},
			expectError: false,
		},
		{
			name:  "Valid URL with HTTP",
			input: "http://example.com",
			expected: &URLParseResult{
				OriginalURL: "http://example.com",
				Normalized:  "http://example.com",
				Domain:      "example.com",
				Path:        "",
				Params:      map[string]string{},
				IsValid:     true,
			},
			expectError: false,
		},
		{
			name:  "URL without scheme (auto-adds HTTPS)",
			input: "example.com",
			expected: &URLParseResult{
				OriginalURL: "https://example.com",
				Normalized:  "https://example.com",
				Domain:      "example.com",
				Path:        "",
				Params:      map[string]string{},
				IsValid:     true,
			},
			expectError: false,
		},
		{
			name:  "URL with trailing slash (normalized)",
			input: "https://example.com/",
			expected: &URLParseResult{
				OriginalURL: "https://example.com/",
				Normalized:  "https://example.com",
				Domain:      "example.com",
				Path:        "/",
				Params:      map[string]string{},
				IsValid:     true,
			},
			expectError: false,
		},
		{
			name:        "Empty URL",
			input:       "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid URL format",
			input:       "not a url",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid domain",
			input:       "https://invalid",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.OriginalURL != tt.expected.OriginalURL {
				t.Errorf("OriginalURL: expected %s, got %s", tt.expected.OriginalURL, result.OriginalURL)
			}

			if result.Normalized != tt.expected.Normalized {
				t.Errorf("Normalized: expected %s, got %s", tt.expected.Normalized, result.Normalized)
			}

			if result.Domain != tt.expected.Domain {
				t.Errorf("Domain: expected %s, got %s", tt.expected.Domain, result.Domain)
			}

			if result.Path != tt.expected.Path {
				t.Errorf("Path: expected %s, got %s", tt.expected.Path, result.Path)
			}

			if !reflect.DeepEqual(result.Params, tt.expected.Params) {
				t.Errorf("Params: expected %v, got %v", tt.expected.Params, result.Params)
			}

			if result.IsValid != tt.expected.IsValid {
				t.Errorf("IsValid: expected %v, got %v", tt.expected.IsValid, result.IsValid)
			}
		})
	}
}

func TestURLParserParseLogEntry(t *testing.T) {
	parser := NewURLParser()

	tests := []struct {
		name     string
		logEntry string
		expected map[string]string
	}{
		{
			name:     "Valid log entry",
			logEntry: `[2023-06-10 12:34:56] "GET /abc123 HTTP/1.1" 301 "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" "192.168.1.1" "https://referrer.com"`,
			expected: map[string]string{
				"shortcode":  "abc123",
				"ip":         "192.168.1.1",
				"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
				"timestamp":  "2023-06-10 12:34:56",
			},
		},
		{
			name:     "Missing parts",
			logEntry: `[2023-06-10 12:34:56] "GET /abc123 HTTP/1.1" 301`,
			expected: map[string]string{
				"shortcode": "abc123",
				"timestamp": "2023-06-10 12:34:56",
			},
		},
		{
			name:     "Empty log entry",
			logEntry: "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseLogEntry(tt.logEntry)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("For key %s: expected %s, got %s", key, expectedValue, result[key])
				}
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d keys, got %d keys", len(tt.expected), len(result))
			}
		})
	}
}
