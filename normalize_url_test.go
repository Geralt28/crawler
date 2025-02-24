package main

import (
	"reflect"
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
			expected: "https://blog.boot.dev/path",
		},
		{
			name:     "Remove HTTP scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "http://blog.boot.dev/path",
		},

		// ✅ Handle trailing slashes
		{
			name:     "Remove trailing slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "https://blog.boot.dev/path",
		},
		{
			name:     "Remove double trailing slash",
			inputURL: "https://blog.boot.dev/path//",
			expected: "https://blog.boot.dev/path",
		},

		// ✅ Mixed cases (scheme, trailing slash)
		{
			name:     "HTTP scheme & trailing slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "http://blog.boot.dev/path",
		},
		{
			name:     "HTTPS scheme & double trailing slash",
			inputURL: "https://blog.boot.dev/path//",
			expected: "https://blog.boot.dev/path",
		},

		// ✅ Edge cases
		{
			name:     "Root domain without path",
			inputURL: "https://blog.boot.dev/",
			expected: "https://blog.boot.dev",
		},
		{
			name:     "Root domain with HTTP",
			inputURL: "http://blog.boot.dev/",
			expected: "http://blog.boot.dev",
		},
		{
			name:     "URL with subdomain",
			inputURL: "https://sub.blog.boot.dev/path",
			expected: "https://sub.blog.boot.dev/path",
		},
		{
			name:     "URL with port",
			inputURL: "https://blog.boot.dev:8080/path",
			expected: "https://blog.boot.dev:8080/path",
		},
		{
			name:     "URL with query parameters (should be removed)",
			inputURL: "https://blog.boot.dev/path?query=123",
			expected: "https://blog.boot.dev/path",
		},
		{
			name:     "URL with fragment (should be removed)",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "https://blog.boot.dev/path",
		},

		// ✅ Unexpected input cases
		{
			name:     "URL without scheme",
			inputURL: "blog.boot.dev/path",
			expected: "http://blog.boot.dev/path",
		},
		{
			name:     "Only domain name",
			inputURL: "blog.boot.dev",
			expected: "http://blog.boot.dev",
		},
		{
			name:     "Uppercase characters in URL",
			inputURL: "HTTPS://BLOG.BOOT.DEV/PATH",
			expected: "https://blog.boot.dev/path",
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

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		// ✅ Test absolute and relative URLs
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="/path/one">Boot.dev</a>
					<a href="https://other.com/path/one">Other Site</a>
				</body>
			</html>`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},

		// ✅ Test multiple absolute URLs
		{
			name:     "multiple absolute URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="https://example.com">Example</a>
					<a href="https://google.com">Google</a>
				</body>
			</html>`,
			expected: []string{"https://example.com", "https://google.com"},
		},

		// ✅ Test multiple relative URLs
		{
			name:     "multiple relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="/about">About</a>
					<a href="/contact">Contact</a>
				</body>
			</html>`,
			expected: []string{"https://blog.boot.dev/about", "https://blog.boot.dev/contact"},
		},

		// ✅ Test missing href attributes
		{
			name:     "ignore missing href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a>Missing HREF</a>
					<a href="">Empty HREF</a>
					<a href="https://valid.com">Valid</a>
				</body>
			</html>`,
			expected: []string{"https://valid.com"},
		},

		// ✅ Test malformed HTML
		{
			name:     "malformed HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="/malformed">Valid Link
				</body>
			</html>`,
			expected: []string{"https://blog.boot.dev/malformed"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
