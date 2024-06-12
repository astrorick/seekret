package seekret

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// this is a server configuration object
type ServerConfig struct {
	ConfigPath string

	// database parameters
	DatabaseType    string `yaml:"databaseType"`
	DatabaseConnStr string `yaml:"databaseConnStr"`

	// http server parameters
	HTTPServerPort uint16 `yaml:"httpServerPort"`
}

// printable server config representation
func (sc *ServerConfig) String() string {
	if sc.ConfigPath == "" {
		return fmt.Sprintf("Default Config: {\n\tDatabase Type: %s\n\tDatabase Connection String: \"%s\"\n\tHTTP Server Port: %d\n}", sc.DatabaseType, sc.DatabaseConnStr, sc.HTTPServerPort)
	} else {
		return fmt.Sprintf("From '%s': {\n\tDatabase Type: %s\n\tDatabase Connection String: \"%s\"\n\tHTTP Server Port: %d\n}", sc.ConfigPath, sc.DatabaseType, sc.DatabaseConnStr, sc.HTTPServerPort)
	}
}

func loadServerConfig(configPath string) (*ServerConfig, error) {
	// see if user provided a config file
	if configPath == "" {
		// no config file provided
		return &ServerConfig{
			ConfigPath:      "",
			DatabaseType:    "sqlite3",
			DatabaseConnStr: "seekret.db",
			HTTPServerPort:  3553,
		}, nil
	}

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

	if serverConfig.DatabaseType == "" {
		serverConfig.DatabaseType = "sqlite3"
	}

	if serverConfig.DatabaseConnStr == "" {
		serverConfig.DatabaseConnStr = "seekret.db"
	}

	if serverConfig.HTTPServerPort == 0 {
		serverConfig.HTTPServerPort = 3553
	}

	// return object
	return &serverConfig, nil
}
