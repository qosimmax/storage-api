package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/qosimmax/storage-api/user"
)

func (c *Client) prepareAddServerStmt() error {
	stmt, err := c.DB.Preparex(`INSERT INTO server(id, name, address) VALUES ($1, $2, $3)
		ON CONFLICT(address) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("error preparing add server stmt %w", err)
	}

	c.AddServerStmt = stmt
	return nil

}

func (c *Client) AddServer(ctx context.Context, data user.ServerData) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	serverID, _ := generateUUID()
	_, err := c.AddServerStmt.ExecContext(cctx, serverID, data.Name, data.Address)
	if err != nil {
		return "", err
	}

	return serverID, nil
}

func (c *Client) prepareFindAvailableServersStmt() error {
	stmt, err := c.DB.Preparex(`with srv as (select id, coalesce(sum(sf.part_size), 0) size
             from server
                      left join server_files sf on server.id = sf.server_id
			group by id)
			select server.id, server.name, server.address
			from srv
					 join server on server.id = srv.id
			order by size
			limit $1`)
	if err != nil {
		return fmt.Errorf("error preparing add find servers stmt %w", err)
	}

	c.FindAvailableServersStmt = stmt
	return nil
}

func (c *Client) FindAvailableServers(ctx context.Context, limit int) (servers []user.ServerData, err error) {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var rows []serverRow
	err = c.FindAvailableServersStmt.SelectContext(cctx, &rows, limit)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		servers = append(servers, user.ServerData{
			ID:      row.ID,
			Name:    row.Name,
			Address: row.Address,
		})
	}

	return servers, nil
}

func (c *Client) prepareFindFileLocationServersStmt() error {
	stmt, err := c.DB.Preparex(`select s.id, s.name, s.address from server_files
				left join server s on s.id = server_files.server_id
				where file_id = $1
				order by "order";`)
	if err != nil {
		return fmt.Errorf("error preparing add find file location servers stmt %w", err)
	}

	c.FindFileLocationServersStmt = stmt
	return nil
}

func (c *Client) FindFileLocationServers(ctx context.Context, fileID string) (servers []user.ServerData, err error) {
	cctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var rows []serverRow
	err = c.FindFileLocationServersStmt.SelectContext(cctx, &rows, fileID)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		servers = append(servers, user.ServerData{
			ID:      row.ID,
			Name:    row.Name,
			Address: row.Address,
		})
	}

	return servers, nil

}

type serverRow struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating uuid: %w", err)
	}
	return id.String(), nil
}
