package main

import (
	"context"
	"log"

	"github.com/qosimmax/storage-api/config"
	"github.com/qosimmax/storage-api/server"
)

func main() {

	log.Println("Starting ...")

	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	var s server.Server

	if err := s.Create(ctx, cfg); err != nil {
		panic(err)
	}

	if err := s.Serve(ctx); err != nil {
		panic(err)
	}
}
