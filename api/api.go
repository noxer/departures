package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Network is the base URL of a transport.rest API, scheme included.
type Network string

const (
	// NetworkVBB covers the whole Berlin-Brandenburg region (incl. regional trains).
	NetworkVBB Network = "https://v6.vbb.transport.rest"
	// NetworkBVG covers Berlin's local transit only (U-Bahn, tram, bus, S-Bahn).
	NetworkBVG Network = "https://v6.bvg.transport.rest"
	// NetworkDB covers Deutsche Bahn's nationwide (Germany) network.
	NetworkDB Network = "https://v6.db.transport.rest"
)

// Networks maps the short names accepted on the command line to their Network.
// All of these are derhuerst-hosted hafas deployments that share this client's
// interface (see https://transport.rest/).
var Networks = map[string]Network{
	"vbb": NetworkVBB,
	"bvg": NetworkBVG,
	"db":  NetworkDB,
}

// ResolveNetwork turns a command-line value into a Network. It accepts a known
// short name (see Networks) or a full base URL of any other compatible hafas
// host (e.g. "https://v6.rmv.transport.rest"). The bool reports whether the
// value was understood.
func ResolveNetwork(name string) (Network, bool) {
	if n, ok := Networks[name]; ok {
		return n, true
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return Network(strings.TrimRight(name, "/")), true
	}
	return "", false
}

type Client struct {
	c *http.Client
	n Network
}

func New(network Network) (*Client, error) {
	return &Client{
		// A timeout keeps a stalled upstream from hanging the CLI forever and
		// lets the caller's retry loop actually kick in.
		c: &http.Client{Timeout: 30 * time.Second},
		n: network,
	}, nil
}

func (c *Client) getJSON(ctx context.Context, v any, urlFormat string, values ...any) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(c.n)+fmt.Sprintf(urlFormat, values...), nil)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	// The API answers errors with a JSON body, but decoding it into the caller's
	// type would silently yield a zero value. Surface the status instead.
	if resp.StatusCode >= http.StatusBadRequest {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("api: unexpected status %s: %s", resp.Status, snippet)
	}

	d := json.NewDecoder(resp.Body)
	return d.Decode(v)
}
