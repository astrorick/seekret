package seekret

import (
	"fmt"
	"log"

	"github.com/astrorick/seekret/pkg/version"
)

type Server struct {
	ServerVersion *version.Version
	ServerConfig  *ServerConfig
}

func NewServer(configPath string) *Server {
	// define server version
	serverVersion := &version.Version{
		Major: 0,
		Minor: 7,
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
	fmt.Printf("Seekret Server v%s by Astrorick\n\n", s.ServerVersion)

	// exit
	return nil
}
