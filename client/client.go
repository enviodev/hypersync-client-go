package client

import (
	"context"
	"encoding/json"
	"fmt"
	arrowhs "github.com/enviodev/hypersync-client-go/arrow"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	ctx    context.Context
	opts   options.Node
	client *http.Client
}

func NewClient(ctx context.Context, opts options.Node) (*Client, error) {
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
	}, nil
}

func (c *Client) GetQueryUrlFromNode(node options.Node) string {
	return strings.Join([]string{node.Endpoint, "query"}, "/")
}

func (c *Client) GeUrlFromNodeAndPath(node options.Node, path ...string) string {
	paths := append([]string{node.Endpoint}, path...)
	return strings.Join(paths, "/")
}

func (c *Client) Get(ctx context.Context, query *types.Query) (*types.QueryResponse[[]types.DataResponse], error) {
	return c.GetArrow(ctx, query)
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

func DoArrow[R any, T types.QueryResponse[[]types.DataResponse]](ctx context.Context, c *Client, url string, method string, payload R) (*T, error) {
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

	if resp.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(responseData))
	}

	defer resp.Body.Close()

	arrowReader, err := arrowhs.NewReader[T](resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse the ipc/arrow response while attempting to read")
	}

	return arrowReader.ResponsePtr(), nil
}
