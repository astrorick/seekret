package seekret

import (
	"fmt"
	"log"

	"github.com/astrorick/seekret/pkg/version"
)

type Server struct {
	serverVersion *version.Version
	serverConfig  *Config
}

func NewServer(configPath string) *Server {
	// define server version
	serverVersion := &version.Version{
		Major: 0,
		Minor: 5,
		Patch: 0,
	}

	// read config file
	serverConfig, err := loadServerConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		serverVersion: serverVersion,
		serverConfig:  serverConfig,
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
	fmt.Printf("Seekret Server v%s by Astrorick\n\n", s.serverVersion)

	// exit
	return nil
}
