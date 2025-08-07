package processor

import (
	"os"
	"testing"
	"path/filepath"
)

func TestNewPDFProcessor(t *testing.T) {
	processor := NewPDFProcessor()

	if processor == nil {
		t.Fatal("NewPDFProcessor returned nil")
	}

	if processor.config == nil {
		t.Fatal("PDFProcessor config is nil")
	}
}

func TestMarginConstants(t *testing.T) {
	expectedMarginMM := 5.0
	expectedPointsPerMM := 2.834645669
	expectedMarginPoints := expectedMarginMM * expectedPointsPerMM

	if MarginMM != expectedMarginMM {
		t.Errorf("MarginMM = %f, want %f", MarginMM, expectedMarginMM)
	}

	if PointsPerMM != expectedPointsPerMM {
		t.Errorf("PointsPerMM = %f, want %f", PointsPerMM, expectedPointsPerMM)
	}

	if MarginPoints != expectedMarginPoints {
		t.Errorf("MarginPoints = %f, want %f", MarginPoints, expectedMarginPoints)
	}
}

func TestProcessPDF_NonExistentFile(t *testing.T) {
	processor := NewPDFProcessor()

	err := processor.ProcessPDF("nonexistent.pdf", "output.pdf")
	if err == nil {
		t.Fatal("Expected error for non-existent input file")
	}
}

func TestProcessPDF_InvalidOutputPath(t *testing.T) {
	processor := NewPDFProcessor()

	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.pdf")

	if err := createMinimalPDF(inputFile); err != nil {
		t.Fatalf("Failed to create test PDF: %v", err)
	}

	err := processor.ProcessPDF(inputFile, "/invalid/path/output.pdf")
	if err == nil {
		t.Fatal("Expected error for invalid output path")
	}
}

func createMinimalPDF(filename string) error {
	content := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
>>
endobj

xref
0 4
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
trailer
<<
/Size 4
/Root 1 0 R
>>
startxref
184
%%EOF`

	return os.WriteFile(filename, []byte(content), 0644)
}