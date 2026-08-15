# Fulfillment Manifest CLI

`manifestctl` accepts one JSON order request on standard input and produces a fulfillment manifest on standard output. It is self-contained: inventory is supplied by the process, and the in-memory repository exists to make the planning and transaction boundaries testable.

Example:

```sh
printf '%s\n' '{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":2}]}' | go run ./cmd/manifestctl
```

The application keeps a clean dependency direction:

- `cmd/manifestctl` handles JSON and process exit behavior.
- `internal/service` orchestrates request normalization, stock allocation, and persistence.
- `internal/domain` owns validation and manifest values.
- `internal/store` provides the transaction and persistence boundary.
