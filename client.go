package hypersyncgo

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	arrowhs "github.com/enviodev/hypersync-client-go/arrow"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type Client struct {
	ctx       context.Context
	opts      options.Node
	client    *http.Client
	rpcClient *ethclient.Client
	userAgent string
}

func NewClient(ctx context.Context, opts options.Node) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid node options")
	}

	rpcConn, err := rpc.DialOptions(ctx, opts.RpcEndpoint, rpc.WithHeader("Authorization", "Bearer "+opts.ApiToken))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to RPC client")
	}
	rpcClient := ethclient.NewClient(rpcConn)

	return &Client{
		ctx:  ctx,
		opts: opts,
		client: &http.Client{
			Timeout: 2 * time.Minute,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		rpcClient: rpcClient,
		userAgent: fmt.Sprintf("hscg/%s", version()),
	}, nil
}

func (c *Client) GetRPC() *ethclient.Client {
	return c.rpcClient
}

// retryJitter returns a random duration in [0, max) using crypto/rand for use in retry backoff.
func retryJitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return time.Duration(n.Int64())
}

// setHeaders sets common headers on the request (API token is sent as Authorization: Bearer).
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.opts.ApiToken)
	req.Header.Set("User-Agent", c.userAgent)
}

func (c *Client) GetQueryUrlFromNode(node options.Node) string {
	return strings.Join([]string{node.Endpoint, "query"}, "/")
}

func (c *Client) GeUrlFromNodeAndPath(node options.Node, path ...string) string {
	paths := append([]string{node.Endpoint}, path...)
	return strings.Join(paths, "/")
}

func (c *Client) Stream(ctx context.Context, query *types.Query, opts *options.StreamOptions) (*Stream, error) {
	stream, err := NewStream(ctx, c, query, opts)
	if err != nil {
		return nil, err
	}

	// TODO: Retries?
	go func() {
		if sErr := stream.Subscribe(); sErr != nil {
			stream.QueueError(sErr)
			return
		}
	}()

	return stream, nil
}

func (c *Client) GetArrow(ctx context.Context, query *types.Query) (*types.QueryResponse, error) {
	base := c.opts.RetryBaseMs

	c.opts.RetryBackoffMs = time.Duration(100)
	c.opts.MaxNumRetries = 3

	var lastErr error
	for i := 0; i < c.opts.MaxNumRetries+1; i++ {
		response, err := DoArrow[*types.Query](ctx, c, c.GeUrlFromNodeAndPath(c.opts, "query", "arrow-ipc"), http.MethodPost, query)
		if err == nil {
			return response, nil
		}
		lastErr = err

		baseMs := base * time.Millisecond

		jitter := retryJitter(c.opts.RetryBackoffMs)

		select {
		case <-time.After(baseMs + jitter):
			base = min(base+c.opts.RetryBackoffMs, c.opts.RetryCeilingMs)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, errors.Wrapf(lastErr, "failed to get arrow data after retries: %d", c.opts.MaxNumRetries)
}

func DoQuery[R any, T any](ctx context.Context, c *Client, method string, payload R) (*T, error) {
	nodeUrl := c.GetQueryUrlFromNode(c.opts)

	reqPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal envio payload")
	}

	req, err := http.NewRequestWithContext(ctx, method, nodeUrl, strings.NewReader(string(reqPayload)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(responseData))
	}

	var result T
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return &result, nil
}

func Do[R any, T any](ctx context.Context, c *Client, url string, method string, payload R) (*T, error) {
	reqPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal envio payload")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(reqPayload)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(responseData))
	}

	var result T
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return &result, nil
}

func DoArrow[R any](ctx context.Context, c *Client, url string, method string, payload R) (*types.QueryResponse, error) {
	reqPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal envio payload")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(reqPayload)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(responseData))
	}

	arrowReader, err := arrowhs.NewQueryResponseReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the ipc/arrow response while attempting to read")
	}

	return arrowReader.GetQueryResponse(), nil
}
