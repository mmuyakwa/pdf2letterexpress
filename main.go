package main

import (
	"fmt"
	"os"

	"github.com/yourorg/pdf2letterexpress/internal/cli"
)

const (
	appName    = "Pdf2LetterExpress"
	appVersion = "1.0.0"
	appDesc    = "PDF converter for LetterExpress compatibility - adds 5mm margins"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	rootCmd := cli.NewRootCommand(appName, appVersion, appDesc)
	return rootCmd.Execute()
}