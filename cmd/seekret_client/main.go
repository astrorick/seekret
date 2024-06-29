package main

import (
	"crypto"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/astrorick/seekret/pkg/srp"
	"github.com/astrorick/seekret/pkg/version"
)

type SeekretClient struct {
	Banner    string           // seekret server banner
	Config    *Config          // server config
	SRPParams *srp.Params      // srp params
	Version   *version.Version // app version
}

func (sc *SeekretClient) Start(configFilePath string) error {
	//fmt.Println(ss.Banner) // TODO: new banner
	fmt.Printf("Seekret Client v%s by Astrorick\n", sc.Version.String())
	fmt.Printf("Local datetime is %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func main() {
	// bind and parse command line flags
	var (
		configFilePath string
		displayHelp    bool
	)
	flag.StringVar(&configFilePath, "c", "", "Client configuration file path.")
	flag.BoolVar(&displayHelp, "h", false, "Display this help message.")
	flag.Parse()

	// display help and exit if help flag is set
	if displayHelp {
		flag.Usage()
		return
	}

	//* default seekret client configuration
	seekretClient := &SeekretClient{
		Banner: " __           _             _   \n" +
			"/ _\\ ___  ___| | ___ __ ___| |_ \n" +
			"\\ \\ / _ \\/ _ \\ |/ / '__/ _ \\ __|\n" +
			"_\\ \\  __/  __/   <| | |  __/ |_ \n" +
			"\\__/\\___|\\___|_|\\_\\_|  \\___|\\__|",
		Config: &Config{},
		SRPParams: &srp.Params{
			SaltSize: 32,
			HashFcn:  crypto.SHA512,
		},
		Version: &version.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	// start seekret server
	if err := seekretClient.Start(configFilePath); err != nil {
		log.Fatal(err)
	}
}
