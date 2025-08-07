package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInputFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		setup    func() (string, func())
	}{
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
			setup:    func() (string, func()) { return "", func() {} },
		},
		{
			name:     "non-existent file",
			filename: "nonexistent.pdf",
			wantErr:  true,
			setup:    func() (string, func()) { return "nonexistent.pdf", func() {} },
		},
		{
			name:     "non-PDF file",
			filename: "test.txt",
			wantErr:  true,
			setup: func() (string, func()) {
				tempDir := os.TempDir()
				filename := filepath.Join(tempDir, "test.txt")
				os.WriteFile(filename, []byte("test"), 0644)
				return filename, func() { os.Remove(filename) }
			},
		},
		{
			name:     "valid PDF file",
			filename: "test.pdf",
			wantErr:  false,
			setup: func() (string, func()) {
				tempDir := os.TempDir()
				filename := filepath.Join(tempDir, "test.pdf")
				content := "%PDF-1.4\ntest content"
				os.WriteFile(filename, []byte(content), 0644)
				return filename, func() { os.Remove(filename) }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, cleanup := tt.setup()
			defer cleanup()

			if filename != "" {
				tt.filename = filename
			}

			err := ValidateInputFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInputFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateOutputFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "document.pdf",
			expected: "document - converted.pdf",
		},
		{
			name:     "filename with path",
			input:    "/path/to/document.pdf",
			expected: filepath.Join("/path/to", "document - converted.pdf"),
		},
		{
			name:     "filename with multiple dots",
			input:    "my.document.v1.pdf",
			expected: "my.document.v1 - converted.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateOutputFilename(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateOutputFilename() = %v, want %v", result, tt.expected)
			}
		})
	}
}