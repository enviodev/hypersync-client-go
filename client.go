package hypersyncgo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/enviodev/hypersync-client-go/pkg/options"
	"github.com/enviodev/hypersync-client-go/pkg/utils"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	ctx      context.Context
	opts     options.Options
	nodeOpts options.Node
	client   *http.Client
}

func NewClient(ctx context.Context, networkId utils.NetworkID, opts options.Options) (*Client, error) {
	nodeOpts, found := opts.GetNodeByNetworkId(networkId)
	if !found {
		return nil, errors.Errorf("failure to create hypersync client - network id %s not found", networkId)
	}

	return &Client{
		ctx:      ctx,
		opts:     opts,
		nodeOpts: *nodeOpts,
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
	}, nil
}

func (c *Client) GetQueryUrlFromNode(node options.Node) string {
	return strings.Join([]string{node.Endpoint, "query"}, "/")
}

func (c *Client) GeUrlFromNodeAndPath(node options.Node, path string) string {
	return strings.Join([]string{node.Endpoint, path}, "/")
}

func DoQuery[R any, T any](ctx context.Context, c *Client, method string, networkId utils.NetworkID, payload R) (*T, error) {
	node, nodeFound := c.opts.GetNodeByNetworkId(networkId)
	if !nodeFound {
		return nil, fmt.Errorf("requested envio network not found: networkId: %s", networkId)
	}

	nodeUrl := c.GetQueryUrlFromNode(*node)

	reqPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal envio payload")
	}

	req, err := http.NewRequestWithContext(ctx, method, nodeUrl, strings.NewReader(string(reqPayload)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}

	req.Header.Set("Content-Type", "application/json")

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

	req.Header.Set("Content-Type", "application/json")

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
