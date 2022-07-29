package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Network string

const (
	NetworkVBB Network = "v5.vbb.transport.rest"
)

type Client struct {
	c *http.Client
	n Network
}

func New(network Network) (*Client, error) {
	return &Client{
		c: http.DefaultClient,
		n: network,
	}, nil
}

func (c *Client) getJSON(ctx context.Context, v interface{}, urlFormat string, values ...interface{}) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(c.n)+fmt.Sprintf(urlFormat, values...), nil)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	d := json.NewDecoder(resp.Body)
	return d.Decode(v)
}
