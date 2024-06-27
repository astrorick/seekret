package server

import (
	"fmt"
	"net/http"

	"github.com/astrorick/seekret/internal/database"

	//? remove driver from here
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	HTTPPort uint16
	Database *database.Database
	/*SRPParams *srp.SRPParams
	JWTParams *JWTParams*/
}

// Start starts the http server with the provided configuration.
func (srv *Server) Start() error {
	// prepare routes for http server
	//http.HandleFunc("/api/create-user", srv.CreateUserRequestHandler())

	// listen and serve
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.HTTPPort), nil); err != nil {
		return err
	}

	// exit
	return nil
}
