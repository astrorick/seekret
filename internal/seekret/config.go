package seekret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type ServerConfig struct {
	ConfigPath      string `json:"configPath"`
	DatabaseConnStr string `json:"databaseConnStr"`
	HTTPServerPort  uint16 `json:"httpServerPort"`
}

func loadServerConfig(configPath string) (*ServerConfig, error) {
	// try to read the provided config file
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// read all bytes from file
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	// unmarshal its content in a structure
	var serverConfig ServerConfig
	if err := json.Unmarshal(configBytes, &serverConfig); err != nil {
		return nil, err
	}

	// check data integrity and assign default values
	serverConfig.ConfigPath = configPath

	if serverConfig.DatabaseConnStr == "" {
		return nil, errors.New("failed to load config file: no 'databaseConnStr' field was provided")
	}

	if serverConfig.HTTPServerPort == 0 {
		serverConfig.HTTPServerPort = 3553
	}

	// return object
	return &serverConfig, nil
}
