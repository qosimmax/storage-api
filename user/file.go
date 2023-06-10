package user

import "context"

// FileInfo is a struct containing data about a uploading file
type FileInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// FileInfoAdder is an interface for adding uploading file info
type FileInfoAdder interface {
	AddFileInfo(ctx context.Context, file FileInfo) error
}

// FileInfoGetter is an interface for getting file info
type FileInfoGetter interface {
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)
}

// PartitionFileInfo is a struct containing data about a partition file
type PartitionFileInfo struct {
	FileID   string `json:"file_id"`
	ServerID string `json:"server_id"`
	PartSize int64  `json:"part_size"`
	Order    int    `json:"order"`
}

// PartitionFileInfoAdder is an interface for adding file info which transferred to server
type PartitionFileInfoAdder interface {
	AddPartitionFileInfo(ctx context.Context, partitionFile PartitionFileInfo) error
}
