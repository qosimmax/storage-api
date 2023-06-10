package database

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/qosimmax/storage-api/config"
)

// Client holds the database client and prepared statements.
type Client struct {
	DB                          *sqlx.DB
	AddServerStmt               *sqlx.Stmt
	FindAvailableServersStmt    *sqlx.Stmt
	AddFileInfoStmt             *sqlx.Stmt
	AddPartitionFileInfoStmt    *sqlx.Stmt
	FindFileLocationServersStmt *sqlx.Stmt
	GetFileInfoStmt             *sqlx.Stmt
}

// Init sets up a new database client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseURL,
		config.DatabasePort,
		config.DatabaseDB,
		config.DatabaseOptions,
	)

	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(config.DatabaseMaxConnections)
	db.SetMaxIdleConns(config.DatabaseMaxIdleConnections)

	c.DB = db

	if err := c.prepareAddServerStmt(); err != nil {
		return err
	}

	if err := c.prepareFindAvailableServersStmt(); err != nil {
		return err
	}

	if err := c.prepareAddFileInfoStmt(); err != nil {
		return err
	}

	if err := c.prepareAddPartitionFileInfoStmt(); err != nil {
		return err
	}

	if err := c.prepareGetFileInfoStmt(); err != nil {
		return err
	}

	if err := c.prepareFindFileLocationServersStmt(); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection and statements.
func (c *Client) Close() error {

	if err := c.AddServerStmt.Close(); err != nil {
		return err
	}

	if err := c.AddFileInfoStmt.Close(); err != nil {
		return err
	}

	if err := c.AddPartitionFileInfoStmt.Close(); err != nil {
		return err
	}

	if err := c.FindAvailableServersStmt.Close(); err != nil {
		return err
	}

	if err := c.GetFileInfoStmt.Close(); err != nil {
		return err
	}

	if err := c.FindFileLocationServersStmt.Close(); err != nil {
		return err
	}

	err := c.DB.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	return nil
}
