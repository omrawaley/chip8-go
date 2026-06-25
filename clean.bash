#!/bin/bash
set -e

output_dir="./bin"

if [[ -d "$output_dir" ]]; then
    rm -rf "$output_dir"
    echo "Cleaned output"
else
    echo "Nothing to clean"
fi
