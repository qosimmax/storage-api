package storage_server

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) SendFile(ctx context.Context, srcFileInfo user.SrcFileInfo, dstServer user.ServerData) (int64, error) {
	//var err error
	conn, err := net.DialTimeout(networkType, dstServer.Address, c.Timeout)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = conn.Close()
	}()

	file, err := srcFileInfo.File.Open()
	if err != nil {
		return 0, err
	}

	defer file.Close()

	_, err = file.Seek(srcFileInfo.Offset, 0)
	if err != nil {
		return 0, err
	}

	// Send the file id
	_, _ = conn.Write([]byte(srcFileInfo.ID))
	// Send the file size
	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(srcFileInfo.Size))
	_, _ = conn.Write(sizeBuf)
	// Send file body
	return io.CopyN(conn, file, srcFileInfo.Size)
}
