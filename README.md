[![Tests Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/test.yml)
[![Security Status](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml/badge.svg)](https://github.com/enviodev/hypersync-client-go/actions/workflows/gosec.yml)
[![Coverage Status](https://coveralls.io/repos/github/enviodev/hypersync-client-go/badge.svg?branch=main)](https://coveralls.io/github/enviodev/hypersync-client-go?branch=main)

# HyperSync Go Client

> WIP - reach out in discord if you need this :)

[envio](https://envio.dev)


## Examples

### Get Head Block

```go
package main

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
)

func main() {
	opts := options.Options{
		Blockchains: []options.Node{
			{
				Type:      utils.EthereumNetwork,
				NetworkId: utils.EthereumNetworkID,
				Endpoint:  "https://eth.hypersync.xyz",
				RpcEndpoint: "https://eth.rpc.hypersync.xyz",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hsClient, err := hypersyncgo.NewHyper(ctx, opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, found := hsClient.GetClient(utils.EthereumNetworkID)
	if !found {
		fmt.Printf("failure to discover hypersync client for network: %d \n", utils.EthereumNetworkID)
		return
	}

	height, err := client.GetHeight(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Current network: %s, height: %d\n", utils.EthereumNetworkID, height)
}
```


## LICENSE

copyright goes here...