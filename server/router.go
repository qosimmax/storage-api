package server

import (
	"net/http"

	"github.com/qosimmax/storage-api/server/internal/handler"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const v1API string = "/storage-api/api/v1"

// setupRoutes - the root route function.
func (s *Server) setupRoutes() {
	s.Router.Handle("/metrics", promhttp.Handler()).Name("Metrics")
	s.Router.HandleFunc("/_healthz", handler.Healthz).Methods(http.MethodGet).Name("Health")

	api := s.Router.PathPrefix(v1API).Subrouter()
	api.HandleFunc("/file-upload", handler.TransferFile(s.DB, s.Storage)).Methods(http.MethodPost).Name("FileUpload")
	api.HandleFunc("/add-server", handler.AddServer(s.DB)).Methods(http.MethodPost).Name("AddServer")
}
