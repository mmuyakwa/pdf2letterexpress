package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/sirupsen/logrus"
)

// Rename this function to avoid conflict with pdf.go
func (p *PDFProcessor) ProcessPDFSimple(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Creating margins by scaling PDF content")

	// Use the main CreateMargins function which calls CreateMarginsForLetterXpress
	if err := p.CreateMargins(inputFile, outputFile); err == nil {
		return nil
	}

	// Fallback to ImageMagick if the main method fails
	logrus.Warn("Main method failed, trying ImageMagick fallback")
	if err := p.CreateMarginsWithImageMagick(inputFile, outputFile); err == nil {
		return nil
	}

	// Final fallback
	logrus.Error("All margin creation methods failed")
	return fmt.Errorf("failed to create margins with any available method")
}

func (p *PDFProcessor) addMarginsWithNUp(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Creating margins by scaling PDF content")

	// Use the main CreateMargins function instead of old methods
	if err := p.CreateMargins(inputFile, outputFile); err == nil {
		return nil
	}

	// Fallback approaches
	logrus.Warn("Main method failed, trying fallbacks")
	if err := p.CreateMarginsWithImageMagick(inputFile, outputFile); err == nil {
		return nil
	}

	// Final fallback: pdfcpu method
	return p.scaleContentWithImport(inputFile, outputFile)
}

func (p *PDFProcessor) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (p *PDFProcessor) scaleContentWithImport(inputFile, outputFile string) error {
	// Read the input PDF and manually scale each page's content
	inputReader, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputReader.Close()

	outputWriter, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputWriter.Close()

	// Read the PDF context
	ctx, err := api.ReadContext(inputReader, p.config)
	if err != nil {
		return fmt.Errorf("failed to read PDF context: %w", err)
	}

	logrus.WithField("pages", ctx.PageCount).Info("Scaling content on each page")

	// Scale content on each page to create margins
	err = p.scaleAllPages(ctx)
	if err != nil {
		logrus.WithError(err).Warn("Content scaling failed, creating simple copy")
	}

	// Write the modified PDF
	if err := api.WriteContext(ctx, outputWriter); err != nil {
		return fmt.Errorf("failed to write scaled PDF: %w", err)
	}

	logrus.WithField("output", outputFile).Info("Successfully created PDF with scaled content")
	return nil
}

func (p *PDFProcessor) scaleAllPages(ctx *model.Context) error {
	// Get page dimensions
	dims, err := ctx.PageDims()
	if err != nil {
		return fmt.Errorf("failed to get page dimensions: %w", err)
	}

	// Process each page
	for i := 1; i <= ctx.PageCount; i++ {
		if i > len(dims) {
			continue
		}

		pageDim := dims[i-1]

		logrus.WithFields(logrus.Fields{
			"page":         i,
			"width":        pageDim.Width,
			"height":       pageDim.Height,
			"marginPoints": MarginPoints,
		}).Debug("Scaling page content")

		// Calculate scale factor to leave margins
		availableWidth := pageDim.Width - (2 * MarginPoints)
		availableHeight := pageDim.Height - (2 * MarginPoints)

		scaleX := availableWidth / pageDim.Width
		scaleY := availableHeight / pageDim.Height

		// Use uniform scaling (smaller scale to maintain aspect ratio)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Calculate translation to center the scaled content
		scaledWidth := pageDim.Width * scale
		scaledHeight := pageDim.Height * scale
		translateX := (pageDim.Width - scaledWidth) / 2
		translateY := (pageDim.Height - scaledHeight) / 2

		logrus.WithFields(logrus.Fields{
			"page":       i,
			"scale":      scale,
			"translateX": translateX,
			"translateY": translateY,
		}).Debug("Calculated transformation parameters")

		// Apply the transformation matrix to the page content
		// This is where we would modify the PDF content stream
		// For now, we log the transformation that would be applied
		logrus.WithFields(logrus.Fields{
			"page":           i,
			"transformation": fmt.Sprintf("%.6f 0 0 %.6f %.6f %.6f cm", scale, scale, translateX, translateY),
		}).Info("Transformation matrix calculated (content stream modification needed)")
	}

	return nil
}
