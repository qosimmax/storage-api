package database

import (
	"context"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) GetServers(ctx context.Context, limit int) ([]user.ServerData, error) {
	servers := []user.ServerData{
		{
			ID:      "100001",
			Name:    "SERVER-1",
			Address: "localhost:7070",
		},
		{
			ID:      "100002",
			Name:    "SERVER-2",
			Address: "localhost:7071",
		},
		{
			ID:      "100003",
			Name:    "SERVER-3",
			Address: "localhost:7072",
		},
		{
			ID:      "100004",
			Name:    "SERVER-4",
			Address: "localhost:7073",
		},
		{
			ID:      "100005",
			Name:    "SERVER-5",
			Address: "localhost:7074",
		},
		{
			ID:      "100006",
			Name:    "SERVER-6",
			Address: "localhost:7075",
		},
	}

	return servers, nil
}
