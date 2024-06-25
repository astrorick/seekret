package server

import (
	"fmt"
	"net/http"

	"github.com/astrorick/seekret/internal/config"
	"github.com/astrorick/seekret/internal/database"
	"github.com/astrorick/seekret/pkg/srp"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	Config    *config.ServerConfig
	Database  *database.Database
	SRPParams *srp.SRPParams
	JWTKey    []byte
}

// Start starts the http server with the provided configuration.
func (srv *Server) Start() error {
	// prepare routes for http server
	http.HandleFunc("/api/create-user", srv.CreateUserRequestHandler())

	// listen and serve
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.Config.HTTPServerPort), nil); err != nil {
		return err
	}

	// exit
	return nil
}
