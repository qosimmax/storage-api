package mock

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/qosimmax/storage-api/user"
)

const DummyData = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. 
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

var dummyServers = []user.ServerData{
	{
		ID:      "1",
		Name:    "MOCK-SERVER-1",
		Address: "1",
	},
	{
		ID:      "2",
		Name:    "MOCK-SERVER-2",
		Address: "2",
	},
	{
		ID:      "3",
		Name:    "MOCK-SERVER-3",
		Address: "3",
	},
	{
		ID:      "4",
		Name:    "MOCK-SERVER-4",
		Address: "4",
	},
	{
		ID:      "5",
		Name:    "MOCK-SERVER-5",
		Address: "5",
	},
	{
		ID:      "6",
		Name:    "MOCK-SERVER-6",
		Address: "6",
	},
}

// DBMock is a struct implements DB functions
type DBMock struct {
}

// FileTransferMock is a struct implements FileTransfer interface
type FileTransferMock struct {
}

// FileReceiverMock is a struct implements FileReceiver interface
type FileReceiverMock struct {
}

func (f FileTransferMock) SendFile(ctx context.Context, srcFileInfo user.SrcFileInfo, dstServer user.ServerData) (sent int64, err error) {
	file, err := srcFileInfo.File.Open()
	if err != nil {
		return 0, err
	}

	defer file.Close()

	_, err = file.Seek(srcFileInfo.Offset, 0)
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("./%s_%s", srcFileInfo.ID, dstServer.ID)
	fw, err := os.Create(path)
	if err != nil {
		return
	}
	defer fw.Close()

	// Send file body
	return io.CopyN(fw, file, srcFileInfo.Size)

}

func (f FileReceiverMock) ReceiveFile(ctx context.Context, fileID string, dstServer user.ServerData, w http.ResponseWriter) (recv int64, err error) {
	path := fmt.Sprintf("./%s_%s", fileID, dstServer.ID)
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	return io.Copy(w, file)
}

func (D DBMock) GetFileInfo(ctx context.Context, fileID string) (*user.FileInfo, error) {
	return &user.FileInfo{ID: fileID, Name: "dummy.txt"}, nil
}

func (D DBMock) FindFileLocationServers(ctx context.Context, fileID string) ([]user.ServerData, error) {
	return dummyServers, nil
}

func (D DBMock) FindAvailableServers(ctx context.Context, limit int) ([]user.ServerData, error) {
	return dummyServers, nil
}

func (D DBMock) AddFileInfo(ctx context.Context, file user.FileInfo) error {
	return nil
}

func (D DBMock) AddPartitionFileInfo(ctx context.Context, partitionFile user.PartitionFileInfo) error {
	return nil
}
