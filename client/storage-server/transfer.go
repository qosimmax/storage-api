package storage_server

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) TransferFile(ctx context.Context, fileInfo user.FileInfo, serverData user.ServerData) (int64, error) {
	var err error
	c.Conns[serverData.ID], err = net.Dial(networkType, serverData.Address)
	if err != nil {
		return 0, err
	}

	file, err := fileInfo.File.Open()
	if err != nil {
		return 0, err
	}

	defer file.Close()

	_, err = file.Seek(fileInfo.Offset, 0)
	if err != nil {
		return 0, err
	}

	// Send the file id
	_, _ = c.Conns[serverData.ID].Write([]byte(fileInfo.ID))
	// Send the file size
	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(fileInfo.Size))
	_, _ = c.Conns[serverData.ID].Write(sizeBuf)
	// Send file body
	return io.CopyN(c.Conns[serverData.ID], file, fileInfo.Size)
}
