package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/astrorick/seekret/internal/config"
	"github.com/astrorick/seekret/pkg/version"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	Version  *version.Version
	Config   *config.ServerConfig
	Database *sql.DB
}

func New(configPath string) (*Server, error) {
	// define server version
	serverVersion := &version.Version{
		Major: 0,
		Minor: 10,
		Patch: 0,
	}

	// read config file
	serverConfig, err := config.NewServerConfig(configPath)
	if err != nil {
		return nil, err
	}

	// when using a 'sqlite3' database, the database file must be created if it does not exist
	if serverConfig.DatabaseType == "sqlite3" {
		if _, err := os.Stat(serverConfig.DatabaseConnStr); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(serverConfig.DatabaseConnStr); err != nil {
				return nil, err
			}
		}
	}

	// open database connection
	serverDB, err := sql.Open(serverConfig.DatabaseType, serverConfig.DatabaseConnStr)
	if err != nil {
		return nil, err
	}

	// run preliminary consistency checks on the server database
	/*if err := srv.runPreliminaryChecks(); err != nil {
		return err
	}*/

	return &Server{
		Config:   serverConfig,
		Database: serverDB,
		Version:  serverVersion,
	}, nil
}

func (srv *Server) Start() error {
	// define banner
	serverBanner := " __           _             _   \n" +
		"/ _\\ ___  ___| | ___ __ ___| |_ \n" +
		"\\ \\ / _ \\/ _ \\ |/ / '__/ _ \\ __|\n" +
		"_\\ \\  __/  __/   <| | |  __/ |_ \n" +
		"\\__/\\___|\\___|_|\\_\\_|  \\___|\\__|"

	// print banner and welcome string
	fmt.Println(serverBanner)
	fmt.Printf("Seekret Server v%s by Astrorick\n", srv.Version)
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	if srv.Config.FilePath != "" {
		fmt.Printf("Using server config at '%s' with parameters:\n\tDatabase Type: %s\n\tDatabase Connecion String: %s\n\tHTTP Server Port: %d\n", srv.Config.FilePath, srv.Config.DatabaseType, srv.Config.DatabaseConnStr, srv.Config.HTTPServerPort)
	} else {
		fmt.Printf("Using the default server config:\n\tDatabase Type: %s\n\tDatabase Connecion String: %s\n\tHTTP Server Port: %d\n", srv.Config.DatabaseType, srv.Config.DatabaseConnStr, srv.Config.HTTPServerPort)
	}

	// enumerate users in database
	var userCount int
	if err := srv.Database.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return err
	}
	fmt.Printf("Found %d registered users\n", userCount)
	defer srv.Database.Close()

	// prepare routes for http server
	http.HandleFunc("/api/create-user", srv.CreateUserRequestHandler())

	// listen and serve
	fmt.Printf("HTTP server listening on port %d\n", srv.Config.HTTPServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", srv.Config.HTTPServerPort), nil); err != nil {
		return err
	}

	// exit
	return nil
}
