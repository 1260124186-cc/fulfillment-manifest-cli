# Docker Validation

Build an architecture-specific image:

```sh
./build_benzhi_docker.sh fulfillment-manifest-cli:arm64 linux/arm64
```

Run the project and send an order request:

```sh
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}' \
  | docker run -i --rm fulfillment-manifest-cli:arm64
```

The image intentionally retains the Go toolchain so validators can run `go build ./...` inside the started container.
