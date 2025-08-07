#!/bin/bash

# Pdf2LetterExpress - Cross-Platform Build Script
# Builds binaries for Windows, Linux, macOS with different architectures

set -e

BINARY_NAME="pdf2letterexpress"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-s -w"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Pdf2LetterExpress Cross-Platform Build${NC}"
echo -e "${BLUE}======================================${NC}"
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Build Time: ${GREEN}${BUILD_TIME}${NC}"
echo ""

# Create bin directory
mkdir -p bin

# Define build targets: GOOS_GOARCH combinations
TARGETS=(
    "windows_amd64:${BINARY_NAME}-windows-amd64.exe"
    "windows_386:${BINARY_NAME}-windows-386.exe"
    "windows_arm64:${BINARY_NAME}-windows-arm64.exe"
    "linux_amd64:${BINARY_NAME}-linux-amd64"
    "linux_386:${BINARY_NAME}-linux-386"
    "linux_arm64:${BINARY_NAME}-linux-arm64"
    "linux_arm:${BINARY_NAME}-linux-arm"
    "darwin_amd64:${BINARY_NAME}-darwin-amd64"
    "darwin_arm64:${BINARY_NAME}-darwin-arm64"
    "freebsd_amd64:${BINARY_NAME}-freebsd-amd64"
    "freebsd_386:${BINARY_NAME}-freebsd-386"
)

# Function to build for specific target
build_target() {
    local target_info=$1
    local goos_goarch=${target_info%:*}
    local output=${target_info#*:}
    local goos=${goos_goarch%_*}
    local goarch=${goos_goarch#*_}
    
    echo -e "${YELLOW}Building for ${goos}/${goarch}...${NC}"
    
    # Set environment variables and build
    env GOOS=${goos} GOARCH=${goarch} CGO_ENABLED=0 \
        go build -ldflags "${LDFLAGS}" -o "bin/${output}" main.go
    
    if [ $? -eq 0 ]; then
        # Get file size
        size=$(ls -lh "bin/${output}" | awk '{print $5}')
        echo -e "${GREEN}✅ Built bin/${output} (${size})${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to build for ${goos}/${goarch}${NC}"
        return 1
    fi
}

# Build all targets
successful_builds=0
total_builds=${#TARGETS[@]}

echo -e "${BLUE}Building ${total_builds} targets...${NC}"
echo ""

for target_info in "${TARGETS[@]}"; do
    if build_target "$target_info"; then
        ((successful_builds++))
    fi
done

echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "${GREEN}Build Summary:${NC}"
echo -e "Successful builds: ${GREEN}${successful_builds}/${total_builds}${NC}"

if [ ${successful_builds} -eq ${total_builds} ]; then
    echo -e "${GREEN}✅ All builds completed successfully!${NC}"
else
    echo -e "${YELLOW}⚠️  Some builds failed.${NC}"
fi

echo ""
echo -e "${BLUE}Generated binaries:${NC}"
ls -la bin/ | grep -E "${BINARY_NAME}" | while read -r line; do
    size=$(echo "$line" | awk '{print $5}')
    name=$(echo "$line" | awk '{print $9}')
    echo -e "  ${GREEN}${name}${NC} (${size} bytes)"
done

echo ""
echo -e "${BLUE}To test a binary:${NC}"
echo -e "  ${YELLOW}./bin/${BINARY_NAME}-linux-amd64 --help${NC}"
echo ""

# Create checksums
echo -e "${BLUE}Generating checksums...${NC}"
cd bin
sha256sum ${BINARY_NAME}-* > checksums.txt 2>/dev/null || shasum -a 256 ${BINARY_NAME}-* > checksums.txt
echo -e "${GREEN}✅ Checksums saved to bin/checksums.txt${NC}"
cd ..

echo -e "${GREEN}Build process completed!${NC}"