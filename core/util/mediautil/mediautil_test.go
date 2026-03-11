package mediautil

import "testing"

func TestIsVideo(t *testing.T) {
	tests := []struct {
		mime string
		want bool
	}{
		{"video/mp4", true},
		{"video/quicktime", true},
		{"audio/mp4", false},
		{"image/jpeg", false},
		{"video", false},
		{"video/mp4/extra", false},
		{"text/plain", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			if got := IsVideo(tt.mime); got != tt.want {
				t.Errorf("IsVideo(%q) = %v, want %v", tt.mime, got, tt.want)
			}
		})
	}
}

func TestIsAudio(t *testing.T) {
	tests := []struct {
		mime string
		want bool
	}{
		{"audio/mpeg", true},
		{"audio/ogg", true},
		{"audio/mp4", true},
		{"video/mp4", false},
		{"image/png", false},
		{"audio", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			if got := IsAudio(tt.mime); got != tt.want {
				t.Errorf("IsAudio(%q) = %v, want %v", tt.mime, got, tt.want)
			}
		})
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		mime string
		want bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"video/mp4", false},
		{"audio/mpeg", false},
		{"image", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			if got := IsImage(tt.mime); got != tt.want {
				t.Errorf("IsImage(%q) = %v, want %v", tt.mime, got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		mime    string
		wantPri string
		wantSub string
		wantOk  bool
	}{
		{"video/mp4", "video", "mp4", true},
		{"text/plain", "text", "plain", true},
		{"video", "", "", false},
		{"video/mp4/extra", "", "", false},
		{"/", "", "", true},
		{"", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			pri, sub, ok := split(tt.mime)
			if pri != tt.wantPri || sub != tt.wantSub || ok != tt.wantOk {
				t.Errorf("split(%q) = (%q, %q, %v), want (%q, %q, %v)", tt.mime, pri, sub, ok, tt.wantPri, tt.wantSub, tt.wantOk)
			}
		})
	}
}
