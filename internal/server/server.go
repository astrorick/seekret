package server

import (
	"fmt"
	"net/http"

	"github.com/astrorick/seekret/internal/database"
	"github.com/astrorick/seekret/pkg/srp"

	//? remove driver from here
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	Database  *database.Database
	Port      uint16
	SRPParams *srp.SRPParams
	JWTKey    []byte
}

// Start starts the http server with the provided configuration.
func (srv *Server) Start() error {
	// prepare routes for http server
	//http.HandleFunc("/api/create-user", srv.CreateUserRequestHandler())

	// listen and serve
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.Port), nil); err != nil {
		return err
	}

	// exit
	return nil
}
