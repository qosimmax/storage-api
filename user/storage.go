package user

import (
	"context"
	"mime/multipart"
	"net/http"
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

// FileReceiver is an interface receive file partition from server
type FileReceiver interface {
	ReceiveFile(ctx context.Context, fileID string, dstServer ServerData, w http.ResponseWriter) (recv int64, err error)
}
