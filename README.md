# hypersync-client-go

[![Tests Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml) [![Security Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml) [![Coverage Status](https://coveralls.io/repos/github/enviodev/hypersync-client-go/badge.svg?branch=main)](https://coveralls.io/github/enviodev/hypersync-client-go?branch=main) [![Discord](https://img.shields.io/badge/Discord-Join%20Chat-7289da?logo=discord&logoColor=white)](https://discord.com/invite/envio)

Go client for [Envio's](https://envio.dev) HyperSync. Provides a native Go interface for accessing HyperSync, HyperRPC, and HyperCURL.

> **Note:** This client is community maintained and marked as work-in-progress. For production use, consider the officially supported clients: [Node.js](https://github.com/enviodev/hypersync-client-node), [Python](https://github.com/enviodev/hypersync-client-python), or [Rust](https://github.com/enviodev/hypersync-client-rust).

## What is HyperSync?

[HyperSync](https://docs.envio.dev/docs/HyperSync/overview) is Envio's high-performance blockchain data retrieval layer. It is a purpose-built alternative to JSON-RPC endpoints, offering up to 2000x faster data access across 70+ EVM-compatible networks and Fuel.

HyperSync lets you query logs, transactions, blocks, and traces with flexible filtering and field selection, returning only the data you need.

If you need a full indexing framework on top of HyperSync with GraphQL APIs and schema management, see [HyperIndex](https://github.com/enviodev/hyperindex).

## Features

- **Native Go interface**: Idiomatic Go API for accessing HyperSync
- **Blocks, logs, transactions, traces**: Query all major blockchain data types
- **HyperRPC support**: Drop-in JSON-RPC compatible endpoint access
- **Streaming**: Process large datasets with built-in pagination
- **Event decoding**: Decode ERC-721 and other ABI-encoded events
- **70+ networks**: Access any [HyperSync-supported network](https://docs.envio.dev/docs/HyperSync/hypersync-supported-networks)

## Installation

```bash
go get github.com/enviodev/hypersync-client-go
```

## API Token

An API token is required to use HyperSync. [Get your token here](https://docs.envio.dev/docs/HyperSync/api-tokens), then set it as an environment variable:

```bash
export ENVIO_API_TOKEN="your-token-here"
```

## Quick Start

```go
import (
    hypersyncgo "github.com/enviodev/hypersync-client-go"
    "github.com/enviodev/hypersync-client-go/options"
)

node := options.Node{
    Endpoint:    "https://eth.hypersync.xyz",
    RpcEndpoint: "https://eth.rpc.hypersync.xyz",
    ApiToken:    os.Getenv("ENVIO_API_TOKEN"),
}

client, err := hypersyncgo.NewClient(ctx, node)
if err != nil {
    log.Fatal(err)
}
```

See the [examples directory](./examples) for complete usage including block ranges, log queries, transaction queries, trace queries, and decoded ERC-721 events.

## Connecting to Different Networks

Change the `Endpoint` to connect to any supported network:

```go
// Arbitrum
node := options.Node{
    Endpoint: "https://arbitrum.hypersync.xyz",
    ApiToken: os.Getenv("ENVIO_API_TOKEN"),
}

// Base
node := options.Node{
    Endpoint: "https://base.hypersync.xyz",
    ApiToken: os.Getenv("ENVIO_API_TOKEN"),
}
```

See the full list of [supported networks and URLs](https://docs.envio.dev/docs/HyperSync/hypersync-supported-networks).

## Running Examples

```bash
# Set your API token first
export ENVIO_API_TOKEN="your-token-here"

# Run an example
go run ./examples/blocks_in_range.go
go run ./examples/logs_in_range.go
go run ./examples/transactions_in_range.go
go run ./examples/traces_in_range.go
go run ./examples/erc721_events_decoded.go
```

Example files use `//go:build ignore` so they are not built with the main module.

## What you can build

The Go client is a good fit for teams building blockchain tooling in Go that need fast, direct access to on-chain data:

- **Blockchain data pipelines**: Query logs, transactions, blocks, and traces at scale without RPC bottlenecks
- **Block explorers**: Power responsive interfaces with real-time data access across 70+ chains
- **Monitoring tools**: Track wallet activity and contract events in near real-time
- **Analytics tools**: Scan entire chain histories in seconds instead of hours
- **Go-native integrations**: Access HyperSync and HyperRPC from existing Go services without switching language

## Testing

```bash
# Load API token from .env file (recommended for local dev)
set -a && source .env && set +a

# Run all tests
go test ./... -v
```

Unit tests that do not call the real API work without a token. Tests that call HyperSync will return 401 without a valid token.

## Documentation

- [HyperSync Documentation](https://docs.envio.dev/docs/HyperSync/overview)
- [HyperRPC Documentation](https://docs.envio.dev/docs/HyperSync/overview-hyperrpc)
- [All Client Libraries](https://docs.envio.dev/docs/HyperSync/hypersync-clients) (Node.js, Python, Rust)
- [Query Reference](https://docs.envio.dev/docs/HyperSync/hypersync-query)
- [Supported Networks](https://docs.envio.dev/docs/HyperSync/hypersync-supported-networks)

## FAQ

**How does this compare to using standard Go JSON-RPC clients?**
HyperSync retrieves data up to 2000x faster than traditional JSON-RPC. It is designed for bulk historical data access across multiple block ranges.

**Do I need an API token?**
Yes. The token must be set on each `options.Node` as `ApiToken`. [Get one here](https://docs.envio.dev/docs/HyperSync/api-tokens).

**Which networks are supported?**
70+ EVM-compatible networks and Fuel. See the [full list](https://docs.envio.dev/docs/HyperSync/hypersync-supported-networks).

**Is this production ready?**
This client is community maintained and marked as work-in-progress. Test thoroughly before using in production. The officially supported clients are [Node.js](https://github.com/enviodev/hypersync-client-node), [Python](https://github.com/enviodev/hypersync-client-python), and [Rust](https://github.com/enviodev/hypersync-client-rust).

**What is HyperRPC?**
HyperRPC is Envio's fast read-only RPC endpoint designed for data-intensive tasks. It is a drop-in JSON-RPC compatible alternative to traditional nodes. Behind the scenes, HyperRPC uses HyperSync to fulfill requests. For most use cases Envio recommends using HyperSync directly, as it provides greater speed and query flexibility. [HyperRPC docs →](https://docs.envio.dev/docs/HyperSync/overview-hyperrpc)

**What is the difference between HyperSync and HyperIndex?**
HyperSync is the raw data access layer. Use it when you need direct access to blockchain data in your own pipeline. [HyperIndex](https://github.com/enviodev/hyperindex) is the full indexing framework built on top of HyperSync, with schema management, event handlers, and a GraphQL API.

## Support

- [Discord community](https://discord.com/invite/envio)
- [GitHub Issues](https://github.com/enviodev/hypersync-client-go/issues)
- [Documentation](https://docs.envio.dev/docs/HyperSync/overview)
