package seekret

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	ConfigPath      string `yaml:"configPath"`
	DatabaseType    string `yaml:"databaseType"`
	DatabaseConnStr string `yaml:"databaseConnStr"`
	HTTPServerPort  uint16 `yaml:"httpServerPort"`
}

// printable config representation
func (sc *ServerConfig) String() string {
	return fmt.Sprintf("%s: {\n\tDatabase Type: %s\n\tDatabase Connection String: \"%s\"\n\tHTTP Server Port: %d\n}", sc.ConfigPath, sc.DatabaseType, sc.DatabaseConnStr, sc.HTTPServerPort)
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
	if err := yaml.Unmarshal(configBytes, &serverConfig); err != nil {
		return nil, err
	}

	// check data integrity and assign default values
	serverConfig.ConfigPath = configPath

	if serverConfig.DatabaseConnStr == "" {
		return nil, errors.New("missing required config field: databaseType")
	}

	if serverConfig.DatabaseConnStr == "" {
		return nil, errors.New("missing required config field: databaseConnStr")
	}

	if serverConfig.HTTPServerPort == 0 {
		serverConfig.HTTPServerPort = 3553
	}

	// return object
	return &serverConfig, nil
}
