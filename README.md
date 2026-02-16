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
    ApiToken:    os.Getenv("HYPERSYNC_API_TOKEN"), // required
}
client, err := hypersyncgo.NewClient(ctx, node)
```

For examples and CI, use the `HYPERSYNC_API_TOKEN` environment variable (and set it in GitHub Actions secrets for integration tests).

## Examples

See more at [Examples](./examples) directory. Examples read the API token from `HYPERSYNC_API_TOKEN`.


## LICENSE

copyright goes here...