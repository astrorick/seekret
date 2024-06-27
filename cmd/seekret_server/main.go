package main

import (
	"crypto"
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/astrorick/seekret/internal/database"
	"github.com/astrorick/seekret/internal/server"
	"github.com/astrorick/seekret/pkg/srp"
	"github.com/astrorick/seekret/pkg/version"

	_ "github.com/mattn/go-sqlite3"
)

type SeekretServer struct {
	Banner    string           // seekret server banner
	Config    *Config          // server config
	SRPParams *srp.SRPParams   // srp params
	Version   *version.Version // app version
}

func (ss *SeekretServer) Start(configFilePath string) error {
	// print banner and welcome string
	//fmt.Println(ss.Banner)
	fmt.Printf("Seekret Server v%s by Astrorick\n", ss.Version.String())
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// load config file from path
	if configFilePath != "" {
		fmt.Printf("Loading server config from '%s'...", configFilePath)

		// build config from config file
		serverConfig, err := NewServerConfig(configFilePath)
		if err != nil {
			return err
		}

		// replace values provided by the server config
		valueDefault := reflect.ValueOf(ss.Config).Elem()
		valueParsed := reflect.ValueOf(serverConfig).Elem()
		for i := 0; i < valueDefault.NumField(); i++ {
			fieldDefault := valueDefault.Field(i)
			fieldParsed := valueParsed.Field(i)

			if fieldParsed.Interface() != reflect.Zero(fieldParsed.Type()).Interface() {
				fieldDefault.Set(fieldParsed)
			}
		}

		fmt.Println("done!")
	} else {
		fmt.Println("Using the default server config.")
	}

	// open connection to database
	fmt.Printf("Connecting to %s database '%s'...", ss.Config.DatabaseType, ss.Config.DatabaseConnStr)
	serverDatabase, err := database.Open(ss.Config.DatabaseType, ss.Config.DatabaseConnStr, ss.Version)
	if err != nil {
		return err
	}
	defer func() {
		if err := serverDatabase.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("done!")

	// enumerate users in database
	userCount, err := serverDatabase.UserCount()
	if err != nil {
		return err
	}
	fmt.Printf("Database contains %d registered user(s)\n", userCount)

	// use the provided (or the default) server config to generate the server object
	srv := &server.Server{
		HTTPPort:  ss.Config.HTTPServerPort,
		Database:  serverDatabase,
		SRPParams: ss.SRPParams,
	}

	// start the http server
	fmt.Printf("Starting HTTP server, listening on port %d", ss.Config.HTTPServerPort)
	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

func main() {
	// bind and parse command line flags
	var (
		configFilePath string
		displayHelp    bool
	)
	flag.StringVar(&configFilePath, "c", "", "Server configuration file path.")
	flag.BoolVar(&displayHelp, "h", false, "Display this help message.")
	flag.Parse()

	// display help and exit if help flag is set
	if displayHelp {
		flag.Usage()
		return
	}

	//* default seekret server configuration
	seekretServer := &SeekretServer{
		Banner: " __           _             _   \n" +
			"/ _\\ ___  ___| | ___ __ ___| |_ \n" +
			"\\ \\ / _ \\/ _ \\ |/ / '__/ _ \\ __|\n" +
			"_\\ \\  __/  __/   <| | |  __/ |_ \n" +
			"\\__/\\___|\\___|_|\\_\\_|  \\___|\\__|",
		Config: &Config{
			// database parameters
			DatabaseType:    "sqlite3",
			DatabaseConnStr: "../../data/seekret.db",

			// http server parameters
			HTTPServerPort: 3553,
		},
		SRPParams: &srp.SRPParams{
			SaltSize: 32,
			HashFcn:  crypto.MD4,
		},
		Version: &version.Version{
			Major: 0,
			Minor: 25,
			Patch: 0,
		},
	}

	// start seekret server
	if err := seekretServer.Start(configFilePath); err != nil {
		log.Fatal(err)
	}
}
