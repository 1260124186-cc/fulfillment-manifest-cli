#!/usr/bin/env sh
set -eu

image_name="${1:-fulfillment-manifest-cli:local}"
platform="${2:-linux/arm64}"

docker build --platform "$platform" -f benzhi.Dockerfile -t "$image_name" .
docker run --rm --platform "$platform" "$image_name" go build ./...
printf '%s\n' '{"order_id":"container-order","customer":"Ada","packages":[{"sku":"book","units":1}]}' \
  | docker run --rm -i --platform "$platform" "$image_name"
