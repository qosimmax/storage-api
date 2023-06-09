package storage_server

import (
	"context"
	"time"

	"github.com/qosimmax/storage-api/config"
)

const networkType = "tcp"

// Client holds the File Server client.
type Client struct {
	Timeout time.Duration
}

// Init sets up a new File Server client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	c.Timeout = 10 * time.Second
	return nil
}
