#!/bin/bash
set -e

app_name="chip8-go"
version="V0.1.1"

src_path="./cmd/tui"
output_dir="./bin"

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

if [[ ! -d "$output_dir" ]]; then
    mkdir -p "$output_dir"
fi

for platform in "${platforms[@]}"; do
    IFS="/" read -r -a pair <<< "$platform"
    GOOS="${pair[0]}"
    GOARCH="${pair[1]}"

    file_ext=""
    if [[ "$GOOS" == "windows" ]]; then
        file_ext=".exe"
    fi

    output_name="${output_dir}/${version}/${GOOS}/${GOARCH}/${app_name}${file_ext}"

    echo "Building $output_name"

    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$output_name" "$src_path"
done

echo "Finished compiling $app_name!"
