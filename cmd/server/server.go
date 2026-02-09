package main

import (
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/httpserver"

	"github.com/rs/zerolog"
)

type Server struct {
	Config config.Config
	Http   httpserver.Server
	Logger zerolog.Logger
	// Services *services.Services
}

type ServerParams struct {
	Config config.Config
	Http   httpserver.Server
	Logger zerolog.Logger
	// Services *services.Services
}

func NewServer(p ServerParams) *Server {
	return &Server{
		Config: p.Config,
		Http:   p.Http,
		Logger: p.Logger,
		// Services: p.Services,
	}
}
