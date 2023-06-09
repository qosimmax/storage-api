package database

import (
	"context"
	"fmt"
	"time"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) prepareAddFileInfoStmt() error {
	stmt, err := c.DB.Preparex(`INSERT INTO file_info(id, name, size) VALUES ($1, $2, $3)`)
	if err != nil {
		return fmt.Errorf("error preparing add file info stmt %w", err)
	}

	c.AddFileInfoStmt = stmt
	return nil
}

func (c *Client) AddFileInfo(ctx context.Context, file user.FileInfo) error {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := c.AddFileInfoStmt.ExecContext(cctx, file.ID, file.Name, file.Size)
	return err
}

func (c *Client) prepareAddPartitionFileInfoStmt() error {
	stmt, err := c.DB.Preparex(`INSERT INTO server_files(file_id, server_id, part_size, "order")
	VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("error preparing add file partition info stmt %w", err)
	}

	c.AddPartitionFileInfoStmt = stmt
	return nil
}

func (c *Client) AddPartitionFileInfo(ctx context.Context, partitionFile user.PartitionFileInfo) error {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := c.AddPartitionFileInfoStmt.ExecContext(cctx, partitionFile.FileID,
		partitionFile.ServerID,
		partitionFile.PartSize,
		partitionFile.Order)

	return err
}

func (c *Client) prepareGetFileInfoStmt() error {
	stmt, err := c.DB.Preparex(`select id, name, size from  file_info where id=$1`)
	if err != nil {
		return fmt.Errorf("error preparing get file info stmt %w", err)
	}

	c.GetFileInfoStmt = stmt
	return nil

}

func (c *Client) GetFileInfo(ctx context.Context, fileID string) (*user.FileInfo, error) {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var row fileInfoRow
	err := c.GetFileInfoStmt.GetContext(cctx, &row, fileID)
	if err != nil {
		return nil, err
	}

	return &user.FileInfo{
		ID:   row.ID,
		Name: row.Name,
		Size: row.Size,
	}, nil

}

type fileInfoRow struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Size int64  `db:"size"`
}
