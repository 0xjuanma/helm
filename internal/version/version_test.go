package version

import (
	"testing"
)

func TestIsOlder(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
		desc string
	}{
		{"v0.12.0", "v0.13.0", true, "older major"},
		{"v0.13.0", "v0.12.0", false, "newer major"},
		{"v0.13.0", "v0.13.0", false, "same version"},
		{"0.12.0", "0.13.0", true, "without v prefix"},
		{"0.13.0", "0.12.0", false, "without v prefix"},
		{"v1.0.0", "v1.0.1", true, "patch version"},
		{"v1.1.0", "v1.0.9", false, "minor version"},
	}

	for _, tt := range tests {
		got := IsOlder(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("IsOlder(%q, %q) = %v; want %v - %s", tt.a, tt.b, got, tt.want, tt.desc)
		}
	}
}
