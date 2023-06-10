package user

import "context"

// ServerData is a struct containing a server info
type ServerData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// ServerFinder is an interface for getting servers which should to transfer files
type ServerFinder interface {
	FindAvailableServers(ctx context.Context, limit int) ([]ServerData, error)
}

// ServerAdder is an interface for adding new server
type ServerAdder interface {
	AddServer(ctx context.Context, data ServerData) (string, error)
}

// FileLocationServerFinder is an interface for find servers which file partitions located
type FileLocationServerFinder interface {
	FindFileLocationServers(ctx context.Context, fileID string) ([]ServerData, error)
}

type ServerFileHandler interface {
	ServerFinder
	FileInfoAdder
	PartitionFileInfoAdder
}
