package storage_server

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) ReceiveFile(ctx context.Context, fileID string, dstServer user.ServerData, w http.ResponseWriter) (int64, error) {
	//var err error
	conn, err := net.DialTimeout(networkType, dstServer.Address, c.Timeout)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = conn.Close()
	}()

	// SEND command
	_, _ = conn.Write([]byte("RECV"))
	// Send the file id
	_, _ = conn.Write([]byte(fileID))
	// receive file body
	return io.Copy(w, conn)
}
