[![Tests Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml)
[![Security Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml)
[![Coverage Status](https://coveralls.io/repos/github/enviodev/hypersync-client-go/badge.svg?branch=main)](https://coveralls.io/github/enviodev/hypersync-client-go?branch=main)

# HyperSync Go Client

> WIP - reach out in discord if you need this :)

Golang client for Envio's HyperSync, HyperRPC and HyperCURL clients.

[Documentation Page](https://docs.envio.dev/docs/hypersync-clients)
<br />
[envio](https://envio.dev)

## Installation

```bash
go get github.com/enviodev/hypersync-client-go
```

## Authentication

The HyperSync API requires an **API token** for all requests. You must set `ApiToken` on each `options.Node` when creating a client. Token and endpoint are validated at client creation; missing or empty token will return an error.

```go
node := options.Node{
    Endpoint:    "https://eth.hypersync.xyz",
    RpcEndpoint: "https://eth.rpc.hypersync.xyz",
    ApiToken:    os.Getenv("ENVIO_API_TOKEN"), // required
}
client, err := hypersyncgo.NewClient(ctx, node)
```

For examples and CI, use the `ENVIO_API_TOKEN` environment variable (and set it in GitHub Actions secrets for integration tests).

## Examples

See more at [Examples](./examples) directory. Examples read the API token from `ENVIO_API_TOKEN`.

## Testing and running examples

### 1. Set your API token

Examples and tests that call the real API need `ENVIO_API_TOKEN` set.

**Option A – use a `.env` file (recommended for local dev)**

In the repo root, create or edit `.env` with an uncommented line:

```bash
ENVIO_API_TOKEN=your-token-here
```

Then load it in your shell before running commands (from the repo root):

```bash
set -a && source .env && set +a
```

**Option B – export in the shell**

```bash
export ENVIO_API_TOKEN="your-token-here"
```

### 2. Run tests

From the repo root:

```bash
go test ./... -v
```

With `.env` loaded:

```bash
set -a && source .env && set +a && go test ./... -v
```

Unit tests (e.g. validation, auth headers) use a placeholder token if `ENVIO_API_TOKEN` is unset. Tests that call the real HyperSync API will 401 unless the token is set.

### 3. Run an example

From the repo root, with the token set (e.g. after `source .env`):

```bash
go run ./examples/blocks_in_range.go
```

Other examples:

- `go run ./examples/logs_in_range.go`
- `go run ./examples/traces_in_range.go`
- `go run ./examples/transactions_in_range.go`

```

Example files use `//go:build ignore` so they are not built with the main module. Running them with `go run ./examples/<name>.go` still compiles and runs that file.

## LICENSE

copyright goes here...
```
