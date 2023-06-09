package storage_server

import (
	"context"
	"net"
	"time"

	"github.com/qosimmax/storage-api/config"
)

const networkType = "tcp"

// Client holds the File Server client.
type Client struct {
	Timeout time.Duration
	Conns   map[string]net.Conn
}

// Init sets up a new File Server client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	c.Conns = make(map[string]net.Conn)
	return nil
}
