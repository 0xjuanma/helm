package cmd

import "testing"

func TestValidateQuickMinutes(t *testing.T) {
	tests := []struct {
		minutes int
		wantErr bool
	}{
		{1, false},
		{60, false},
		{10, false},
		{0, true},
		{61, true},
		{-5, true},
	}

	for _, tt := range tests {
		err := validateQuickMinutes(tt.minutes)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateQuickMinutes(%d) error = %v, wantErr %v", tt.minutes, err, tt.wantErr)
		}
	}
}
