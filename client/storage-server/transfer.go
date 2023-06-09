package storage_server

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) TransferFile(ctx context.Context, fileInfo user.FileInfo, serverData user.ServerData) (int64, error) {
	//var err error
	conn, err := net.DialTimeout(networkType, serverData.Address, c.Timeout)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = conn.Close()
	}()

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
	_, _ = conn.Write([]byte(fileInfo.ID))
	// Send the file size
	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(fileInfo.Size))
	_, _ = conn.Write(sizeBuf)
	// Send file body
	return io.CopyN(conn, file, fileInfo.Size)
}
