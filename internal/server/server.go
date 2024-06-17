package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/astrorick/seekret/internal/config"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	Config   *config.ServerConfig
	Database *sql.DB
}

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
