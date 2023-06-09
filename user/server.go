package user

import "context"

// ServerData is a struct containing a server info
type ServerData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// ServerFinder is an interface for getting server info
type ServerFinder interface {
	GetServers(ctx context.Context, limit int) ([]ServerData, error)
}

// ServerAdder is an interface for adding new server
type ServerAdder interface {
	AddServer(ctx context.Context, data ServerData) error
}
