#!/bin/bash

# PROJECT_VERSION_FILE_PATH=$1
# if [ -z "$PROJECT_VERSION_FILE_PATH" ]; then
#     echo "PROJECT_VERSION_FILE_PATH is required"
#     exit 1
# fi

PROJECT_VERSION_FILE_PATH="version/version.go"

if [ -z "$PROJECT_VERSION" ]; then
    read -rp "PROJECT_VERSION:v" PROJECT_VERSION
fi

if [ -z "$PROJECT_VERSION" ]; then
    echo "PROJECT_VERSION is required"
    exit 1
fi

echo "PROJECT_VERSION: $PROJECT_VERSION"

# modify version.go (Version = "x.x.x")
sedr 'Version[\s]+=[\s]+\"(.*?)\"' "\$1" "$PROJECT_VERSION" "$PROJECT_VERSION_FILE_PATH"
