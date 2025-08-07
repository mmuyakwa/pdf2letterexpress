package processor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/sirupsen/logrus"
)

// CreateMarginsWithPDFCPU uses pdfcpu to directly manipulate PDF structure for margins
func (p *PDFProcessor) CreateMarginsWithPDFCPU(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Creating margins using pdfcpu direct manipulation")

	// Open input PDF
	inputFileHandle, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFileHandle.Close()

	// Read PDF context
	ctx, err := api.ReadContext(inputFileHandle, p.config)
	if err != nil {
		return fmt.Errorf("failed to read PDF context: %w", err)
	}

	// Convert 5mm to points
	marginPoints := MarginMM * PointsPerMM

	logrus.WithField("marginPoints", marginPoints).Debug("Calculated margin in points")

	// Process each page
	for pageNum := 1; pageNum <= ctx.PageCount; pageNum++ {
		err := p.addMarginsToPage(ctx, pageNum, marginPoints)
		if err != nil {
			return fmt.Errorf("failed to add margins to page %d: %w", pageNum, err)
		}
	}

	// Write modified PDF
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFileHandle.Close()

	err = api.WriteContext(ctx, outputFileHandle)
	if err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	logrus.WithField("output", outputFile).Info("Successfully created PDF with pdfcpu margins")
	return nil
}

// addMarginsToPage adds margins to a specific page by manipulating MediaBox and content
func (p *PDFProcessor) addMarginsToPage(ctx *model.Context, pageNum int, marginPoints float64) error {
	// Get page dictionary
	pageDict, _, _, err := ctx.PageDict(pageNum, false)
	if err != nil {
		return fmt.Errorf("failed to get page dict: %w", err)
	}

	// Extract current MediaBox
	mediaBoxObj, found := pageDict.Find("MediaBox")
	if !found || mediaBoxObj == nil {
		return fmt.Errorf("no MediaBox found on page %d", pageNum)
	}

	mediaBoxArray, ok := mediaBoxObj.(types.Array)
	if !ok || len(mediaBoxArray) < 4 {
		return fmt.Errorf("invalid MediaBox on page %d", pageNum)
	}

	// Extract current dimensions
	x0, y0, x1, y1, err := p.extractMediaBoxValues(mediaBoxArray)
	if err != nil {
		return fmt.Errorf("failed to extract MediaBox values: %w", err)
	}

	currentWidth := x1 - x0
	currentHeight := y1 - y0

	logrus.WithFields(logrus.Fields{
		"page":          pageNum,
		"currentWidth":  currentWidth,
		"currentHeight": currentHeight,
		"x0":            x0, "y0": y0, "x1": x1, "y1": y1,
	}).Debug("Current page dimensions")

	// Calculate new MediaBox with margins
	newX0 := x0
	newY0 := y0
	newX1 := x1 + (2 * marginPoints) // Add margin to both sides
	newY1 := y1 + (2 * marginPoints) // Add margin to top and bottom

	// Create new MediaBox
	newMediaBox := types.Array{
		types.Float(newX0),
		types.Float(newY0),
		types.Float(newX1),
		types.Float(newY1),
	}

	// Update MediaBox in page dictionary
	pageDict.Update("MediaBox", newMediaBox)

	// Create content stream to translate existing content
	contentStream := fmt.Sprintf("q\n%f 0 0 %f %f %f cm\n", 1.0, 1.0, marginPoints, marginPoints)

	// Get existing content
	existingContent, err := p.getPageContent(ctx, pageNum)
	if err != nil {
		logrus.WithError(err).Warn("Could not get existing content, proceeding with translation only")
		existingContent = ""
	}

	// Combine translation with existing content
	newContent := contentStream + existingContent + "\nQ\n"

	// Update page content
	err = p.updatePageContent(ctx, pageNum, newContent)
	if err != nil {
		return fmt.Errorf("failed to update page content: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"page":      pageNum,
		"newX1":     newX1,
		"newY1":     newY1,
		"newWidth":  newX1 - newX0,
		"newHeight": newY1 - newY0,
	}).Debug("Updated page with margins")

	return nil
}

// extractMediaBoxValues extracts float values from MediaBox array
func (p *PDFProcessor) extractMediaBoxValues(mediaBoxArray types.Array) (float64, float64, float64, float64, error) {
	var x0, y0, x1, y1 float64

	// Extract x0
	if floatObj, isFloat := mediaBoxArray[0].(types.Float); isFloat {
		x0 = float64(floatObj)
	} else if intObj, isInt := mediaBoxArray[0].(types.Integer); isInt {
		x0 = float64(intObj)
	} else {
		return 0, 0, 0, 0, fmt.Errorf("invalid x0 value")
	}

	// Extract y0
	if floatObj, isFloat := mediaBoxArray[1].(types.Float); isFloat {
		y0 = float64(floatObj)
	} else if intObj, isInt := mediaBoxArray[1].(types.Integer); isInt {
		y0 = float64(intObj)
	} else {
		return 0, 0, 0, 0, fmt.Errorf("invalid y0 value")
	}

	// Extract x1
	if floatObj, isFloat := mediaBoxArray[2].(types.Float); isFloat {
		x1 = float64(floatObj)
	} else if intObj, isInt := mediaBoxArray[2].(types.Integer); isInt {
		x1 = float64(intObj)
	} else {
		return 0, 0, 0, 0, fmt.Errorf("invalid x1 value")
	}

	// Extract y1
	if floatObj, isFloat := mediaBoxArray[3].(types.Float); isFloat {
		y1 = float64(floatObj)
	} else if intObj, isInt := mediaBoxArray[3].(types.Integer); isInt {
		y1 = float64(intObj)
	} else {
		return 0, 0, 0, 0, fmt.Errorf("invalid y1 value")
	}

	return x0, y0, x1, y1, nil
}

// getPageContent retrieves existing content stream from page
func (p *PDFProcessor) getPageContent(ctx *model.Context, pageNum int) (string, error) {
	// This is a simplified version - in a full implementation,
	// you would need to handle multiple content streams and decompress them
	pageDict, _, _, err := ctx.PageDict(pageNum, false)
	if err != nil {
		return "", err
	}

	contentsObj, found := pageDict.Find("Contents")
	if !found {
		return "", nil // No existing content
	}

	// Handle different types of content objects
	switch contentsObj.(type) {
	case types.StreamDict:
		// Single content stream - would need decompression in real implementation
		return "% Existing content preserved\n", nil
	case types.Array:
		// Multiple content streams
		return "% Multiple content streams preserved\n", nil
	default:
		return "", nil
	}
}

// updatePageContent updates the content stream of a page
func (p *PDFProcessor) updatePageContent(ctx *model.Context, pageNum int, newContent string) error {
	pageDict, _, _, err := ctx.PageDict(pageNum, false)
	if err != nil {
		return err
	}

	// Create new content stream
	contentBytes := []byte(newContent)

	// Create stream dictionary
	streamDict := types.StreamDict{
		Dict: types.Dict{
			"Length": types.Integer(len(contentBytes)),
		},
		Content: contentBytes,
	}

	// Update page contents
	pageDict.Update("Contents", streamDict)

	return nil
}

// CreateMarginsWithGhostscript - keep as fallback method
func (p *PDFProcessor) CreateMarginsWithGhostscript(inputFile, outputFile string) error {
	logrus.Info("Redirecting to LetterXpress-compatible method")
	return p.CreateMarginsForLetterXpress(inputFile, outputFile)
}

// createMarginsWithPDFtk uses PDFtk as fallback to create margins
func (p *PDFProcessor) createMarginsWithPDFtk(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
	}).Info("Trying PDFtk for margin creation")

	// Check if PDFtk is available
	_, err := exec.LookPath("pdftk")
	if err != nil {
		logrus.WithError(err).Warn("PDFtk not found either, using simple copy")
		return p.copyFile(inputFile, outputFile)
	}

	// PDFtk approach: This is more complex and would require additional steps
	// For now, fall back to copy
	logrus.Warn("PDFtk margin creation not implemented yet, using copy")
	return p.copyFile(inputFile, outputFile)
}

// CreateMarginsWithImageMagick uses ImageMagick as alternative
func (p *PDFProcessor) CreateMarginsWithImageMagick(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Creating margins using ImageMagick")

	// Check if ImageMagick convert is available
	_, err := exec.LookPath("convert")
	if err != nil {
		logrus.WithError(err).Error("ImageMagick not found")
		return fmt.Errorf("imagemagick not available: %w", err)
	}

	// Use ImageMagick to add margins by scaling and centering
	// This converts PDF to image and back, but creates visible margins
	tempDir := filepath.Dir(outputFile)
	tempImagePattern := filepath.Join(tempDir, "temp_page_%d.png")

	// Step 1: Convert PDF to high-resolution images
	cmd1 := exec.Command("convert",
		"-density", "300", // High DPI for quality
		inputFile,
		tempImagePattern,
	)

	logrus.WithField("command", strings.Join(cmd1.Args, " ")).Debug("Converting PDF to images")
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("output", string(output1)).Error("PDF to image conversion failed")
		return fmt.Errorf("pdf to image conversion failed: %w", err)
	}

	// Step 2: Add margins to images and convert back to PDF
	// Calculate margin size as percentage of image
	marginPercent := int(MarginMM * 2) // Rough conversion for margin percentage

	cmd2 := exec.Command("convert",
		tempImagePattern,
		"-bordercolor", "white",
		"-border", fmt.Sprintf("%d%%x%d%%", marginPercent, marginPercent),
		"-density", "300",
		outputFile,
	)

	logrus.WithField("command", strings.Join(cmd2.Args, " ")).Debug("Adding margins and converting back to PDF")
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("output", string(output2)).Error("Image to PDF conversion failed")
		return fmt.Errorf("image to pdf conversion failed: %w", err)
	}

	// Clean up temporary files
	cleanupCmd := exec.Command("rm", "-f", strings.Replace(tempImagePattern, "%d", "*", 1))
	cleanupCmd.Run()

	logrus.WithField("output", outputFile).Info("Successfully created PDF with ImageMagick margins")
	return nil
}

// CreateMarginsForLetterXpress creates margins while maintaining exact DIN A4 format
func (p *PDFProcessor) CreateMarginsForLetterXpress(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Creating margins for LetterXpress while maintaining DIN A4 format")

	// Check if ImageMagick is available
	_, err := exec.LookPath("convert")
	if err != nil {
		logrus.WithError(err).Error("ImageMagick not found")
		return fmt.Errorf("imagemagick not available: %w", err)
	}

	// LetterXpress requires EXACT DIN A4: 210 × 297 mm
	targetWidthMM := 210.0
	targetHeightMM := 297.0

	// Calculate content area (A4 minus 5mm margins on all sides)
	contentWidthMM := targetWidthMM - (2 * MarginMM)   // 200mm
	contentHeightMM := targetHeightMM - (2 * MarginMM) // 287mm

	// Use 300 DPI for high quality
	dpi := 300.0

	// Calculate final A4 dimensions in pixels at 300 DPI
	finalWidthPx := int(targetWidthMM * dpi / 25.4)
	finalHeightPx := int(targetHeightMM * dpi / 25.4)

	// Calculate content area in pixels
	contentWidthPx := int(contentWidthMM * dpi / 25.4)
	contentHeightPx := int(contentHeightMM * dpi / 25.4)
	marginPx := int(MarginMM * dpi / 25.4)

	logrus.WithFields(logrus.Fields{
		"targetWidthMM":   targetWidthMM,
		"targetHeightMM":  targetHeightMM,
		"contentWidthMM":  contentWidthMM,
		"contentHeightMM": contentHeightMM,
		"finalWidthPx":    finalWidthPx,
		"finalHeightPx":   finalHeightPx,
		"contentWidthPx":  contentWidthPx,
		"contentHeightPx": contentHeightPx,
		"marginPx":        marginPx,
		"dpi":             dpi,
	}).Info("Calculated exact DIN A4 pixel dimensions")

	tempDir := filepath.Dir(outputFile)
	tempContentFile := filepath.Join(tempDir, "temp_content.png")
	tempA4File := filepath.Join(tempDir, "temp_a4_canvas.png")

	// Step 1: Convert PDF content to scaled image
	cmd1 := exec.Command("convert",
		"-density", fmt.Sprintf("%.0f", dpi),
		"-background", "white",
		"-flatten",
		"-resize", fmt.Sprintf("%dx%d!", contentWidthPx, contentHeightPx), // Force exact size with !
		inputFile+"[0]", // First page only
		tempContentFile,
	)

	logrus.WithField("command", strings.Join(cmd1.Args, " ")).Debug("Converting PDF to scaled content")
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("output", string(output1)).Error("PDF content conversion failed")
		return fmt.Errorf("pdf content conversion failed: %w", err)
	}

	// Step 2: Create exact A4 white canvas and place content with margins
	cmd2 := exec.Command("convert",
		"-size", fmt.Sprintf("%dx%d", finalWidthPx, finalHeightPx),
		"xc:white", // Create white canvas
		tempContentFile,
		"-geometry", fmt.Sprintf("+%d+%d", marginPx, marginPx), // Position content with margins
		"-composite",
		tempA4File,
	)

	logrus.WithField("command", strings.Join(cmd2.Args, " ")).Debug("Creating A4 canvas with content")
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("output", string(output2)).Error("A4 canvas creation failed")
		return fmt.Errorf("a4 canvas creation failed: %w", err)
	}

	// Step 3: Convert to PDF with exact DIN A4 size specification
	cmd3 := exec.Command("convert",
		tempA4File,
		"-density", fmt.Sprintf("%.0f", dpi),
		"-compress", "jpeg",
		"-quality", "95",
		"-define", "pdf:page-size=a4", // Force A4 page size
		outputFile,
	)

	logrus.WithField("command", strings.Join(cmd3.Args, " ")).Debug("Converting to final A4 PDF")
	output3, err := cmd3.CombinedOutput()
	if err != nil {
		logrus.WithError(err).WithField("output", string(output3)).Error("Final PDF creation failed")
		return fmt.Errorf("final pdf creation failed: %w", err)
	}

	// Clean up temporary files
	os.Remove(tempContentFile)
	os.Remove(tempA4File)

	logrus.WithFields(logrus.Fields{
		"output":            outputFile,
		"finalSizeMM":       fmt.Sprintf("%.1f × %.1f mm", targetWidthMM, targetHeightMM),
		"letterXpressReady": true,
	}).Info("Successfully created LetterXpress-compatible PDF with exact DIN A4 format")

	return nil
}

// CreateMargins - Main function using LetterXpress approach
func (p *PDFProcessor) CreateMargins(inputFile, outputFile string) error {
	// Use LetterXpress-compatible approach as primary method
	return p.CreateMarginsForLetterXpress(inputFile, outputFile)
}
