package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateInputFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	if !isPDFFile(filename) {
		return fmt.Errorf("file must be a PDF: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	n, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("cannot read file header: %w", err)
	}

	if n < 4 || !strings.HasPrefix(string(buffer), "%PDF") {
		return fmt.Errorf("file is not a valid PDF: missing PDF header")
	}

	return nil
}

func isPDFFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".pdf"
}

func ValidateOutputPath(filename string) error {
	dir := filepath.Dir(filename)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", dir)
	}

	tempFile := filepath.Join(dir, ".pdf2letterexpress_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("output directory is not writable: %w", err)
	}

	file.Close()
	os.Remove(tempFile)

	return nil
}