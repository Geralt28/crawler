package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		// ✅ Basic cases
		{
			name:     "Remove HTTPS scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Remove HTTP scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},

		// ✅ Handle trailing slashes
		{
			name:     "Remove trailing slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Remove double trailing slash",
			inputURL: "https://blog.boot.dev/path//",
			expected: "blog.boot.dev/path",
		},

		// ✅ Mixed cases (scheme, trailing slash)
		{
			name:     "HTTP scheme & trailing slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "HTTPS scheme & double trailing slash",
			inputURL: "https://blog.boot.dev/path//",
			expected: "blog.boot.dev/path",
		},

		// ✅ Edge cases
		{
			name:     "Root domain without path",
			inputURL: "https://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
		{
			name:     "Root domain with HTTP",
			inputURL: "http://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
		{
			name:     "URL with subdomain",
			inputURL: "https://sub.blog.boot.dev/path",
			expected: "sub.blog.boot.dev/path",
		},
		{
			name:     "URL with port",
			inputURL: "https://blog.boot.dev:8080/path",
			expected: "blog.boot.dev:8080/path",
		},
		{
			name:     "URL with query parameters (should be removed)",
			inputURL: "https://blog.boot.dev/path?query=123",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "URL with fragment (should be removed)",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path",
		},

		// ✅ Unexpected input cases
		{
			name:     "URL without scheme",
			inputURL: "blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "Only domain name",
			inputURL: "blog.boot.dev",
			expected: "blog.boot.dev",
		},
		{
			name:     "Uppercase characters in URL",
			inputURL: "HTTPS://BLOG.BOOT.DEV/PATH",
			expected: "blog.boot.dev/path",
		},

		// ✅ Handle empty input
		{
			name:     "Empty string input",
			inputURL: "",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
