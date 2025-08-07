# Usage Guide

## Quick Start

The most basic usage requires only the PDF file path:

```bash
pdf2letterexpress document.pdf
```

This will create `document - converted.pdf` in the same directory.

## Command Syntax

```
pdf2letterexpress [flags] <PDF-file>
```

### Flags

| Flag          | Short | Description                              | Default |
| ------------- | ----- | ---------------------------------------- | ------- |
| `--verbose`   | `-v`  | Enable verbose logging                   | `false` |
| `--log-level` |       | Set log level (debug, info, warn, error) | `info`  |
| `--version`   |       | Show version information                 |         |
| `--help`      | `-h`  | Show help message                        |         |

## Examples

### Basic Conversion

```bash
pdf2letterexpress invoice.pdf
```

**Output**: `invoice - converted.pdf`

### Verbose Mode

```bash
pdf2letterexpress -v report.pdf
```

Shows detailed processing information.

### Debug Mode

```bash
pdf2letterexpress --log-level debug presentation.pdf
```

Shows maximum detail for troubleshooting.

### Paths with Spaces

```bash
pdf2letterexpress "My Document.pdf"
pdf2letterexpress '/path/with spaces/document.pdf'
```

## Output File Naming

The output file is always created in the same directory as the input file with the suffix " - converted.pdf":

- `document.pdf` → `document - converted.pdf`
- `report-2024.pdf` → `report-2024 - converted.pdf`
- `invoice.v2.pdf` → `invoice.v2 - converted.pdf`

## Error Messages

### Common Errors

**File not found**:

```
Error: input validation failed: file does not exist: document.pdf
```

**Not a PDF file**:

```
Error: input validation failed: file must be a PDF: document.txt
```

**Permission denied**:

```
Error: failed to create output file: permission denied
```

**Corrupted PDF**:

```
Error: PDF processing failed: failed to read PDF context
```

## Processing Details

### What the Tool Does

1. **Validates** the input PDF file
1. **Analyzes** page dimensions
1. **Calculates** scaling factors to create 5mm margins
1. **Processes** each page with content scaling
1. **Writes** the converted PDF to the output file

### Margin Calculation

- **Target margin**: 5mm on all sides
- **Conversion**: 5mm = ~14.17 points (PDF units)
- **Scaling**: Content scaled to fit within reduced page area
- **Centering**: Scaled content centered on the page

## Best Practices

### File Organization

```bash
# Organize files in folders
mkdir converted
pdf2letterexpress document.pdf
mv "document - converted.pdf" converted/
```

### Batch Processing

```bash
# Process multiple files
for file in *.pdf; do
    pdf2letterexpress "$file"
done
```

### Quality Verification

```bash
# Use verbose mode to verify processing
pdf2letterexpress -v important-document.pdf
```

## Troubleshooting

### High Memory Usage

For very large PDF files, ensure sufficient system memory is available.

### Slow Processing

Large files with many pages will take longer. Use verbose mode to monitor progress.

### Output Quality

The tool preserves all content quality. If output appears scaled, this is the intended 5mm margin effect.