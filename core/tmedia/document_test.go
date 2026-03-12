package tmedia

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetDocumentName(t *testing.T) {
	tests := []struct {
		name string
		doc  *tg.Document
		want string
	}{
		{
			name: "with filename attribute",
			doc: &tg.Document{
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeFilename{FileName: "test.txt"},
				},
			},
			want: "test.txt",
		},
		{
			name: "without filename attribute, known mime type",
			doc: &tg.Document{
				ID:       12345,
				MimeType: "application/pdf",
			},
			want: "12345.pdf",
		},
		{
			name: "without filename attribute, unknown mime type",
			doc: &tg.Document{
				ID:       67890,
				MimeType: "foo/bar",
			},
			want: "67890.unknown",
		},
		{
			name: "without filename attribute, empty mime type",
			doc: &tg.Document{
				ID:       112233,
				MimeType: "",
			},
			want: "112233.unknown",
		},
		{
			name: "with multiple attributes, including filename",
			doc: &tg.Document{
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeAudio{Duration: 10},
					&tg.DocumentAttributeFilename{FileName: "audio.mp3"},
				},
			},
			want: "audio.mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDocumentName(tt.doc); got != tt.want {
				t.Errorf("GetDocumentName() = %v, want %v", got, tt.want)
			}
		})
	}
}
