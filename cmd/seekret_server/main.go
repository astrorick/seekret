package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/astrorick/seekret/internal/config"
	"github.com/astrorick/seekret/internal/database"
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
	flag.StringVar(&configFilePath, "config", "", "Server configuration file path.")
	flag.BoolVar(&displayHelp, "help", false, "Display this help message.")
	flag.Parse()

	// display help and exit if help flag is set
	if displayHelp {
		flag.Usage()
		return
	}

	// app version
	appVersion := &version.Version{
		Major: 0,
		Minor: 21,
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
	fmt.Printf("Seekret Server v%s by Astrorick\n", appVersion.String())
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// load config file
	var (
		serverConfig *config.ServerConfig
		err          error
	)
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

	// open connection to database
	serverDatabase, err := database.Open(serverConfig.DatabaseType, serverConfig.DatabaseConnStr, appVersion)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := serverDatabase.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// enumerate users in database
	userCount, err := serverDatabase.UserCount()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Database contains %d registered user(s)\n", userCount)

	// start http server with provided settings
	srv := &server.Server{
		Config:   serverConfig,
		Database: serverDatabase,
	}
	fmt.Printf("Starting HTTP server on port %d\n", srv.Config.HTTPServerPort)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
