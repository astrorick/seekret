package seekret

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	Config   *ServerConfig
	Database *sql.DB
	Version  *Version
}

func NewServer(configPath string) (*Server, error) {
	// read config file
	serverConfig, err := loadServerConfig(configPath)
	if err != nil {
		return nil, err
	}

	// init database connection
	serverDB, err := sql.Open(serverConfig.DatabaseType, serverConfig.DatabaseConnStr)
	if err != nil {
		return nil, err
	}

	// define server version
	serverVersion := &Version{
		Major: 0,
		Minor: 10,
		Patch: 0,
	}

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

	// enumerate users in database (also checks if database is ok)
	var userCount int
	if err := srv.Database.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return err
	}
	fmt.Printf("Database OK with %d registered users\n", userCount)
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
