package server

import (
	"fmt"
	"net/http"

	"github.com/astrorick/seekret/internal/database"
	"github.com/astrorick/seekret/pkg/jwt"
	"github.com/astrorick/seekret/pkg/srp"
)

type Server struct {
	HTTPPort uint16
	Database *database.Database
	SRP      *srp.Params
	JWT      *jwt.Params
}

// Start starts the http server with the provided configuration.
func (srv *Server) Start() error {
	// prepare routes for http server
	http.HandleFunc("/api/user/create", srv.CreateUserRequestHandler())
	http.HandleFunc("/api/user/delete", srv.DeleteUserRequestHandler())

	// listen and serve
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.HTTPPort), nil); err != nil {
		return err
	}

	// exit
	return nil
}
