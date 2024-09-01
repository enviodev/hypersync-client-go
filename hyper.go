package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/pkg/errors"
	"sync"
)

// Hyper manages a collection of blockchain clients.
type Hyper struct {
	ctx     context.Context
	opts    options.Options
	mu      *sync.RWMutex
	clients map[utils.NetworkID]*Client
}

// NewHyper creates a new instance of HyperSync with the given context and options.
// It validates the provided options and initializes clients for each blockchain network.
//
// Returns an error if the options are invalid or if a client for any network cannot be created.
func NewHyper(ctx context.Context, opts options.Options) (*Hyper, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid options to hyper client")
	}

	zLog, err := logger.GetZapLogger("development", opts.LogLevel.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize zap logger")
	}
	logger.SetGlobalLogger(zLog)

	mu := &sync.RWMutex{}
	clientMap := make(map[utils.NetworkID]*Client)

	for _, clientOpts := range opts.GetBlockchains() {
		mu.Lock()
		if _, ok := clientMap[clientOpts.NetworkId]; !ok {
			nClient, err := NewClient(ctx, clientOpts)
			if err != nil {
				mu.Unlock()
				return nil, errors.Wrapf(err, "failed to create hypersync client for network %s", clientOpts.NetworkId)
			}
			clientMap[clientOpts.NetworkId] = nClient
		}
		mu.Unlock()
	}

	return &Hyper{
		ctx:     ctx,
		opts:    opts,
		mu:      mu,
		clients: clientMap,
	}, nil
}

// GetClients returns a map of all blockchain clients managed by HyperSync.
func (h *Hyper) GetClients() map[utils.NetworkID]*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients
}

// GetClient returns a specific blockchain client by its network ID.
// The boolean return value indicates whether the client was found.
func (h *Hyper) GetClient(networkId utils.NetworkID) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[networkId]
	return c, ok
}
