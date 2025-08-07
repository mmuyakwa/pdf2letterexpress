#!/bin/bash

if ! go build -o pdf2letterexpress main.go; then
    echo "Build failed. Aborting."
    exit 1
fi

sh build.sh

./bin/pdf2letterexpress-darwin-arm64 --help


cp bin/pdf2letterexpress-darwin-arm64 .

chmod +x pdf2letterexpress-darwin-arm64

# Check if a PDF file parameter is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <path-to-pdf-file>"
    echo "Please provide a PDF file to test with."
    exit 1
fi

./pdf2letterexpress-darwin-arm64 "$1"
