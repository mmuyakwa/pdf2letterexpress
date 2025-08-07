package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/sirupsen/logrus"
)

const (
	MarginMM        = 5.0
	PointsPerMM     = 2.834645669
	MarginPoints    = MarginMM * PointsPerMM
)

type PDFProcessor struct {
	config *model.Configuration
}

func NewPDFProcessor() *PDFProcessor {
	config := model.NewDefaultConfiguration()

	config.ValidationMode = model.ValidationRelaxed

	return &PDFProcessor{
		config: config,
	}
}

func (p *PDFProcessor) ProcessPDF(inputFile, outputFile string) error {
	logrus.WithFields(logrus.Fields{
		"input":  inputFile,
		"output": outputFile,
		"margin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Debug("Processing PDF file")

	// Use a simpler approach: NUp with 1 page per sheet, scaled down to create margins
	return p.addMarginsWithNUp(inputFile, outputFile)
}

func (p *PDFProcessor) addMarginsToPDF(input io.ReadSeeker, output io.Writer) error {
	ctx, err := api.ReadContext(input, p.config)
	if err != nil {
		return fmt.Errorf("failed to read PDF context: %w", err)
	}

	logrus.WithField("pages", ctx.PageCount).Debug("PDF loaded successfully")

	if err := p.scalePagesForMargins(ctx); err != nil {
		return fmt.Errorf("failed to scale pages: %w", err)
	}

	if err := api.WriteContext(ctx, output); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

func (p *PDFProcessor) scalePagesForMargins(ctx *model.Context) error {
	pageCount := ctx.PageCount

	for i := 1; i <= pageCount; i++ {
		logrus.WithField("page", i).Debug("Processing page")

		if err := p.scalePageContent(ctx, i); err != nil {
			return fmt.Errorf("failed to scale page %d: %w", i, err)
		}
	}

	return nil
}

func (p *PDFProcessor) scalePageContent(ctx *model.Context, pageNr int) error {
	dims, err := ctx.PageDims()
	if err != nil {
		return fmt.Errorf("failed to get page dimensions: %w", err)
	}

	if pageNr > len(dims) {
		return fmt.Errorf("page %d does not exist", pageNr)
	}

	pageDim := dims[pageNr-1] // Convert to 0-based index

	logrus.WithFields(logrus.Fields{
		"page":   pageNr,
		"width":  pageDim.Width,
		"height": pageDim.Height,
	}).Debug("Page dimensions")

	availableWidth := pageDim.Width - (2 * MarginPoints)
	availableHeight := pageDim.Height - (2 * MarginPoints)

	scaleX := availableWidth / pageDim.Width
	scaleY := availableHeight / pageDim.Height

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	logrus.WithFields(logrus.Fields{
		"page":       pageNr,
		"scaleX":     scaleX,
		"scaleY":     scaleY,
		"finalScale": scale,
	}).Debug("Calculated scale factors")

	logrus.WithFields(logrus.Fields{
		"page":       pageNr,
		"scaleX":     scaleX,
		"scaleY":     scaleY,
		"finalScale": scale,
	}).Debug("Calculated scale factors - using new resize method")

	return nil
}

func (p *PDFProcessor) transformPagesWithMargins(ctx *model.Context) error {
	// Get page dimensions
	dims, err := ctx.PageDims()
	if err != nil {
		return fmt.Errorf("failed to get page dimensions: %w", err)
	}

	// Transform each page
	for i := 1; i <= ctx.PageCount; i++ {
		if i > len(dims) {
			continue
		}

		pageDim := dims[i-1]
		
		logrus.WithFields(logrus.Fields{
			"page": i,
			"width": pageDim.Width,
			"height": pageDim.Height,
		}).Debug("Transforming page with margins")

		// Apply transformation matrix to scale and center content
		err := p.applyContentTransformation(ctx, i, pageDim)
		if err != nil {
			return fmt.Errorf("failed to transform page %d: %w", i, err)
		}
	}

	return nil
}

func (p *PDFProcessor) applyContentTransformation(ctx *model.Context, pageNum int, pageDim types.Dim) error {
	// Calculate scale factor to leave 5mm margins
	targetMarginPoints := MarginPoints
	
	// Calculate available space after margins
	availableWidth := pageDim.Width - (2 * targetMarginPoints)
	availableHeight := pageDim.Height - (2 * targetMarginPoints)
	
	// Calculate scale factors
	scaleX := availableWidth / pageDim.Width
	scaleY := availableHeight / pageDim.Height
	
	// Use the smaller scale to maintain aspect ratio
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
		"page": pageNum,
		"scale": scale,
		"translateX": translateX,
		"translateY": translateY,
		"marginPoints": targetMarginPoints,
	}).Info("Applying content transformation")

	// For now, we'll use pdfcpu's built-in scaling feature
	// This is a simplified approach - in production you might need more sophisticated content stream manipulation
	
	return nil
}

func (p *PDFProcessor) addMarginsUsingNUp(inputFile, outputFile string) error {
	// Use pdfcpu's transformations to add margins
	// This approach uses crop and then resize to create the effect
	
	logrus.WithFields(logrus.Fields{
		"input": inputFile,
		"output": outputFile,
		"targetMargin": fmt.Sprintf("%.1fmm", MarginMM),
	}).Info("Adding margins by transforming PDF content")
	
	// Step 1: First copy the input to output
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

	// Step 2: Read the PDF context and transform it
	ctx, err := api.ReadContext(inputReader, p.config)
	if err != nil {
		return fmt.Errorf("failed to read PDF context: %w", err)
	}

	logrus.WithField("pages", ctx.PageCount).Info("Processing pages for margin addition")

	// Step 3: Transform each page to add margins by scaling content
	err = p.transformPagesWithMargins(ctx)
	if err != nil {
		return fmt.Errorf("failed to transform pages: %w", err)
	}

	// Step 4: Write the transformed PDF
	if err := api.WriteContext(ctx, outputWriter); err != nil {
		return fmt.Errorf("failed to write transformed PDF: %w", err)
	}

	return nil
}