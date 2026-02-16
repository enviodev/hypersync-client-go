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

A **bearer token is required** to use the HyperSync client. The token is sent as an `Authorization: Bearer <token>` header on all HTTP requests to the HyperSync server.

Set the `BearerToken` field when configuring a node:

```go
opts := options.Options{
    Blockchains: []options.Node{
        {
            Type:        utils.EthereumNetwork,
            NetworkId:   utils.EthereumNetworkID,
            Endpoint:    "https://eth.hypersync.xyz",
            RpcEndpoint: "https://eth.rpc.hypersync.xyz",
            BearerToken: os.Getenv("HYPERSYNC_BEARER_TOKEN"),
        },
    },
}
```

Client creation will fail if `BearerToken` is empty.

## Examples

See more at [Examples](./examples) directory.


## LICENSE

copyright goes here...