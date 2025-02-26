package utils

import (
	"testing"
)

func TestURLJoin(t *testing.T) {
	tests := []struct {
		base     string
		elements []string
		expected string
		wantErr  bool
	}{
		{"https://example.com", []string{"path", "to", "resource"}, "https://example.com/path/to/resource", false},
		{"https://example.com/base", []string{"sub"}, "https://example.com/base/sub", false},
		{"https://example.com/", []string{"/leading", "/slashes"}, "https://example.com/leading/slashes", false},
		{"not a url", []string{"invalid"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.base, func(t *testing.T) {
			got, err := URLJoin(tt.base, tt.elements...)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLJoin() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.expected {
				t.Errorf("URLJoin() = %v, want %v", got, tt.expected)
			}
		})
	}
}
