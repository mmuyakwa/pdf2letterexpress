package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GenerateOutputFilename(inputFile string) string {
	dir := filepath.Dir(inputFile)
	base := filepath.Base(inputFile)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	outputName := nameWithoutExt + " - converted.pdf"
	return filepath.Join(dir, outputName)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GetFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func EnsureDirectoryExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}