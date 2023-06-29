#!/bin/bash

# Set the package name and version
if [ ! "$PACKAGE_NAME" ]; then
    echo -ne "\nMissing Package Name.\n"
    exit 1
fi

if [ ! "$PACKAGE_VERSION" ]; then
    echo -ne "\nMissing Package Version.\n"
    exit 1
fi

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
  file_name="${PACKAGE_NAME}-${PACKAGE_VERSION}-${GOOS}-${GOARCH}"
  output_file="$OUTPUT_DIR/$file_name"

  if [ $GOOS = "windows" ]; then
        output_zip="${output_file}"
        output_file+='.exe'
        GOOS=$GOOS GOARCH=$GOARCH go build -o "${output_file}"
        zip -m "${output_zip}.zip" "${output_file}"
    else
        GOOS=$GOOS GOARCH=$GOARCH go build -o "${output_file}"
        tar -czvf "${output_file}.tar.gz" "${output_file}"
    fi

    rm "${output_file}"
done
