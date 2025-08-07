package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"github.com/yourorg/pdf2letterexpress/internal/processor"
	"github.com/yourorg/pdf2letterexpress/internal/utils"
)

type Config struct {
	InputFile  string
	OutputFile string
	Verbose    bool
	LogLevel   string
}

func NewRootCommand(appName, appVersion, appDesc string) *cobra.Command {
	config := &Config{}

	rootCmd := &cobra.Command{
		Use:     fmt.Sprintf("%s <PDF-file>", appName),
		Short:   appDesc,
		Long:    fmt.Sprintf("%s\n\n%s", appDesc, "Automatically scales PDF content to create 5mm margins on all sides for LetterExpress compatibility."),
		Version: appVersion,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConversion(config, args[0])
		},
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&config.LogLevel, "log-level", "info", "Set log level (debug, info, warn, error)")

	return rootCmd
}

func runConversion(config *Config, inputFile string) error {
	setupLogging(config)

	logrus.WithField("input", inputFile).Info("Starting PDF conversion")

	if err := utils.ValidateInputFile(inputFile); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}

	outputFile := utils.GenerateOutputFilename(inputFile)
	config.InputFile = inputFile
	config.OutputFile = outputFile

	logrus.WithField("output", outputFile).Info("Output file will be created")

	processor := processor.NewPDFProcessor()
	if err := processor.ProcessPDF(inputFile, outputFile); err != nil {
		return fmt.Errorf("PDF processing failed: %w", err)
	}

	fmt.Printf("✅ Successfully converted PDF\n")
	fmt.Printf("📁 Input:  %s\n", inputFile)
	fmt.Printf("📁 Output: %s\n", outputFile)

	return nil
}

func setupLogging(config *Config) {
	level := logrus.InfoLevel
	if config.Verbose {
		level = logrus.DebugLevel
	}

	if config.LogLevel != "" {
		parsedLevel, err := logrus.ParseLevel(config.LogLevel)
		if err == nil {
			level = parsedLevel
		}
	}

	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "15:04:05",
		FullTimestamp:   true,
	})
}