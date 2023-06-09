package database

import (
	"context"

	"github.com/qosimmax/storage-api/config"
)

// Client holds the database client and prepared statements.
type Client struct {
}

// Init sets up a new database client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	return nil
}
