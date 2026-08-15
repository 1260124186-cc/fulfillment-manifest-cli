#!/usr/bin/env sh
set -eu

image_name="${1:-fulfillment-manifest-cli:local}"
platform="${2:-linux/arm64}"

docker build --platform "$platform" -f benzhi.Dockerfile -t "$image_name" .
