package seekret

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/astrorick/seekret/internal/api"
	"github.com/astrorick/seekret/pkg/version"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Server struct {
	ServerVersion *version.Version
	ServerConfig  *ServerConfig
}

func NewServer(configPath string) *Server {
	// define server version
	serverVersion := &version.Version{
		Major: 0,
		Minor: 10,
		Patch: 0,
	}

	// read config file
	serverConfig, err := loadServerConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		ServerVersion: serverVersion,
		ServerConfig:  serverConfig,
	}
}

func (s *Server) Start() error {
	// define banner
	serverBanner := " __           _             _   \n" +
		"/ _\\ ___  ___| | ___ __ ___| |_ \n" +
		"\\ \\ / _ \\/ _ \\ |/ / '__/ _ \\ __|\n" +
		"_\\ \\  __/  __/   <| | |  __/ |_ \n" +
		"\\__/\\___|\\___|_|\\_\\_|  \\___|\\__|"

	// print banner and welcome string
	fmt.Println(serverBanner)
	fmt.Printf("Seekret Server v%s by Astrorick\n", s.ServerVersion)
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// initialize database object
	db, err := sql.Open(s.ServerConfig.DatabaseType, s.ServerConfig.DatabaseConnStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// enumerate users in database (also checks if database is ok)
	var userCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return err
	}
	fmt.Printf("Successfully loaded '%s': found %d registered users\n", s.ServerConfig.DatabaseConnStr, userCount)

	// prepare routes for http server
	http.HandleFunc("/api/create-user", api.CreateUserRequestHandler(db))

	// listen and serve
	fmt.Printf("HTTP server listening on port %d\n", s.ServerConfig.HTTPServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.ServerConfig.HTTPServerPort), nil); err != nil {
		return err
	}

	// exit
	return nil
}
