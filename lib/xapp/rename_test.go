package xapp

import (
	"testing"
	"time"
)

func TestRename(t *testing.T) {
	// Fixed time for reproducible tests
	fixedTime := time.Date(2023, 10, 27, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		cfg        Config
		path       string
		wantFormat string // simpler check than exact string due to time components if not careful, but we mocked time
	}{
		{
			name: "Normal Rename",
			cfg: Config{
				Rename: "{year}/{fname}{ext}",
			},
			path:       "test/file.png",
			wantFormat: "2023/file.png",
		},
		{
			name: "HMAC Rename Full",
			cfg: Config{
				Rename:     "{hmac}/{fname}{ext}",
				HmacKey:    "secret_key",
				HmacFormat: "{fname}{ext}",
			},
			path: "test/image.jpg",
			// hmac_sha256("image.jpg", "secret_key") = 9614d443a1f6bd260dc20455b008e3e7a35a4e8bd11951082821dffd27ccb070
			wantFormat: "9614d443a1f6bd260dc20455b008e3e7a35a4e8bd11951082821dffd27ccb070/image.jpg",
		},
		{
			name: "HMAC Rename Truncated",
			cfg: Config{
				Rename:     "{hmac}/{fname}{ext}",
				HmacKey:    "secret_key",
				HmacFormat: "{fname}{ext}",
				HmacLen:    8,
			},
			path: "test/image.jpg",
			// first 8 chars of hash
			wantFormat: "9614d443/image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppCfg = tt.cfg
			got := Rename(tt.path, fixedTime)
			if got != tt.wantFormat {
				t.Errorf("Rename() = %v, want %v", got, tt.wantFormat)
			}
		})
	}
}
