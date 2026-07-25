package cmd

import (
	"errors"
	"testing"
)

func TestDecideUpdate(t *testing.T) {
	tests := []struct {
		name        string
		current     string
		latest      string
		fetchErr    error
		wantProceed bool
	}{
		{"dev build skips update", "dev", "v1.0.0", nil, false},
		{"fetch error still proceeds", "v0.9.0", "", errors.New("network error"), true},
		{"empty latest still proceeds", "v0.9.0", "", nil, true},
		{"already on latest short-circuits", "v1.0.0", "v1.0.0", nil, false},
		{"local ahead of latest short-circuits", "v1.1.0", "v1.0.0", nil, false},
		{"older version proceeds", "v0.9.0", "v1.0.0", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proceed, message := decideUpdate(tt.current, tt.latest, tt.fetchErr)
			if proceed != tt.wantProceed {
				t.Errorf("decideUpdate(%q, %q, %v) proceed = %v; want %v", tt.current, tt.latest, tt.fetchErr, proceed, tt.wantProceed)
			}
			if !proceed && message == "" {
				t.Errorf("decideUpdate(%q, %q, %v) returned no message when short-circuiting", tt.current, tt.latest, tt.fetchErr)
			}
		})
	}
}
