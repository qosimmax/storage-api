package user

import (
	"context"
	"mime/multipart"
)

type FileInfo struct {
	ID     string
	File   *multipart.FileHeader
	Offset int64
	Size   int64
}

type FileTransfer interface {
	TransferFile(ctx context.Context, fileInfo FileInfo, serverData ServerData) (sent int64, err error)
}
