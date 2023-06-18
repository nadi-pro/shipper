#!/bin/bash

# Set the package name and version
PACKAGE_NAME="Nadi-Shipper"
VERSION="1.0.0"

# Define the target operating systems and architectures with file extensions
TARGETS=(
  "linux/amd64:tar.gz"
  "linux/arm64:tar.gz"
  "darwin/amd64:tar.gz"
  "windows/amd64:zip"
)

# Set the output directory
OUTPUT_DIR="installers"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Build for each target
for target in "${TARGETS[@]}"; do
  # Split the target into GOOS and GOARCH with file extension
  goos_arch="${target%:*}"
  file_extension="${target#*:}"
  IFS='/' read -r goos goarch <<< "$goos_arch"

  # Set the environment variables
  export GOOS="$goos"
  export GOARCH="$goarch"

  # Build the package
  echo "Building for $GOOS/$GOARCH..."
  file_name="${PACKAGE_NAME}_${VERSION}_${GOOS}_${GOARCH}"
  file_name_title_case="$(echo "$file_name" | tr '[:lower:]' '[:upper:]')"
  output_file="$OUTPUT_DIR/$file_name_title_case.$file_extension"

  # Build the package with the appropriate file extension
  case "$file_extension" in
    "tar.gz")
      go build -o "$output_file" .
      ;;
    "zip")
      go build -o "$output_file" .
      zip -j "$output_file" "$output_file"
      ;;
    *)
      echo "Unsupported file extension: $file_extension"
      ;;
  esac
done
