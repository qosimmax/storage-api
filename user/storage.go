package user

import (
	"context"
	"mime/multipart"
)

// SrcFileInfo is a struct contains data about a transfer file
type SrcFileInfo struct {
	ID     string
	File   *multipart.FileHeader
	Offset int64
	Size   int64
}

// FileTransfer is an interface for transfer file to destination server
type FileTransfer interface {
	SendFile(ctx context.Context, srcFileInfo SrcFileInfo, dstServer ServerData) (sent int64, err error)
}
