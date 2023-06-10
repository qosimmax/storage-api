// Package server provides functionality to easily set up an HTTTP server.
//
// The server holds all the clients it needs and they should be set up in the Create method.
//
// The HTTP routes and middleware are set up in the setupRouter method.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/qosimmax/storage-api/client/database"
	storage_server "github.com/qosimmax/storage-api/client/storage-server"
	"github.com/qosimmax/storage-api/config"

	"github.com/gorilla/mux"
)

// Server holds the HTTP server, router, config and all clients.
type Server struct {
	Config  *config.Config
	DB      *database.Client
	Storage *storage_server.Client
	HTTP    *http.Server
	Router  *mux.Router
}

// Create sets up the HTTP server, router and all clients.
// Returns an error if an error occurs.
func (s *Server) Create(ctx context.Context, config *config.Config) error {

	var dbClient database.Client
	if err := dbClient.Init(ctx, config); err != nil {
		return fmt.Errorf("database client: %w", err)
	}

	var storageClient storage_server.Client
	if err := storageClient.Init(ctx, config); err != nil {
		return fmt.Errorf("storage client: %w", err)
	}

	s.DB = &dbClient
	s.Storage = &storageClient
	s.Config = config
	s.Router = mux.NewRouter()
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Config.Port),
		Handler: s.Router,
	}

	s.setupRoutes()

	return nil
}

// Serve tells the server to start listening and serve HTTP requests.
// It also makes sure that the server gracefully shuts down on exit.
// Returns an error if an error occurs.
func (srv *Server) Serve(ctx context.Context) error {
	idleConnsClosed := make(chan struct{}) // this is used to signal that we can not exit
	go func(ctx context.Context, s *http.Server) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		<-stop

		log.Println("Shutdown signal received")

		if err := s.Shutdown(ctx); err != nil {
			log.Println(err.Error())
		}

		if err := srv.DB.Close(); err != nil {
			log.Println(err.Error())
		}

		close(idleConnsClosed) // call close to say we can now exit the function
	}(ctx, srv.HTTP)

	log.Printf("Ready at: %s\n", srv.Config.Port)

	if err := srv.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("unexpected server error: %w", err)
	}
	<-idleConnsClosed // this will block until close is called

	return nil
}
