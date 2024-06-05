package seekret

import (
	"fmt"
	"log"

	"github.com/astrorick/seekret/pkg/version"
)

type Seekret struct {
	serverVersion *version.Version
	serverConfig  *Config
}

func New(configPath string) *Seekret {
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

	return &Seekret{
		serverVersion: serverVersion,
		serverConfig:  serverConfig,
	}
}

func (s *Seekret) Start() error {
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
