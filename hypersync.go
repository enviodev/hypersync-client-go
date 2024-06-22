package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/client"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/pkg/errors"
	"sync"
)

type HyperSync struct {
	ctx     context.Context
	opts    options.Options
	mu      sync.RWMutex
	clients map[utils.NetworkID]*client.Client
}

func NewHyperSync(ctx context.Context, opts options.Options) (*HyperSync, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid options to hypersync client")
	}

	mu := sync.RWMutex{}
	clientMap := make(map[utils.NetworkID]*client.Client)

	for _, clientOpts := range opts.GetBlockchains() {
		mu.Lock()
		if _, ok := clientMap[clientOpts.NetworkId]; !ok {
			nClient, err := client.NewClient(ctx, clientOpts)
			if err != nil {
				mu.Unlock()
				return nil, errors.Wrapf(err, "failed to create hypersync client for network %s", clientOpts.NetworkId)
			}
			clientMap[clientOpts.NetworkId] = nClient
		}
		mu.Unlock()
	}

	return &HyperSync{
		ctx:     ctx,
		opts:    opts,
		mu:      mu,
		clients: clientMap,
	}, nil
}

func (h *HyperSync) GetClients() map[utils.NetworkID]*client.Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients
}

func (h *HyperSync) GetClient(networkId utils.NetworkID) (*client.Client, bool) {
	fmt.Println(h.clients)
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[networkId]
	return c, ok
}
