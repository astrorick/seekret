package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/astrorick/seekret/internal/config"
	"github.com/astrorick/seekret/internal/server"
	"github.com/astrorick/seekret/pkg/version"
)

func main() {
	// define flags to match command line arguments
	var (
		configFilePath string
		displayHelp    bool
	)

	// bind and parse command line flags
	flag.StringVar(&configFilePath, "config", "", "Configuration file path.")
	flag.BoolVar(&displayHelp, "help", false, "Display help.")
	flag.Parse()

	// display help and exit if help flag is set
	if displayHelp {
		flag.Usage()
		return
	}

	// app version
	appVersion := &version.Version{
		Major: 0,
		Minor: 20,
		Patch: 0,
	}

	// define banner
	appBanner := " __           _             _   \n" +
		"/ _\\ ___  ___| | ___ __ ___| |_ \n" +
		"\\ \\ / _ \\/ _ \\ |/ / '__/ _ \\ __|\n" +
		"_\\ \\  __/  __/   <| | |  __/ |_ \n" +
		"\\__/\\___|\\___|_|\\_\\_|  \\___|\\__|"

	// print banner and welcome string
	fmt.Println(appBanner)
	fmt.Printf("Seekret Server v%s by Astrorick\n", appVersion)
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// read/load config file
	var serverConfig *config.ServerConfig
	var err error
	if configFilePath != "" {
		// build config from config file
		serverConfig, err = config.NewServerConfig(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(
			"Using server config at '%s' with parameters:\n\tDatabase Type: %s\n\tDatabase Connecion String: %s\n\tHTTP Server Port: %d\n",
			serverConfig.FilePath,
			serverConfig.DatabaseType,
			serverConfig.DatabaseConnStr,
			serverConfig.HTTPServerPort)
	} else {
		serverConfig = config.DefaultServerConfig
		fmt.Printf(
			"Using the default server config:\n\tDatabase Type: %s\n\tDatabase Connecion String: %s\n\tHTTP Server Port: %d\n",
			serverConfig.DatabaseType,
			serverConfig.DatabaseConnStr,
			serverConfig.HTTPServerPort)
	}

	// when using a 'sqlite3' database, the database file must be created if it does not exist
	if serverConfig.DatabaseType == "sqlite3" {
		if _, err := os.Stat(serverConfig.DatabaseConnStr); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(serverConfig.DatabaseConnStr); err != nil {
				log.Fatal(err)
			}
		}
	}

	// open database connection
	serverDB, err := sql.Open(serverConfig.DatabaseType, serverConfig.DatabaseConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: run db consistency checks

	// enumerate users in database
	var userCount uint64
	if err := serverDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Database contains %d registered users\n", userCount)
	defer serverDB.Close()

	// start http server with provided settings
	srv := &server.Server{
		Config:   serverConfig,
		Database: serverDB,
	}
	fmt.Printf("Starting HTTP server on port %d\n", srv.Config.HTTPServerPort)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
