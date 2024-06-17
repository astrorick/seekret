package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	FilePath string

	// database parameters
	DatabaseType    string `yaml:"databaseType"`
	DatabaseConnStr string `yaml:"databaseConnStr"`

	// http server parameters
	HTTPServerPort uint16 `yaml:"httpServerPort"`
}

var DefaultServerConfig = &ServerConfig{
	FilePath:        "",
	DatabaseType:    "sqlite3",
	DatabaseConnStr: "seekret.db",
	HTTPServerPort:  3553,
}

// NewServerConfig tries to open, read and parse the provided file in a 'ServerConfig' object.
// An error is returned if the provided file does not exist or if the file is not accessible.
func NewServerConfig(filePath string) (*ServerConfig, error) {
	// open specified config file
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// read the entire file
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	// parse its content in a struct
	var serverConfig ServerConfig
	if err := yaml.Unmarshal(configBytes, &serverConfig); err != nil {
		return nil, err
	}

	// assign default values when necessary
	if serverConfig.DatabaseType == "" {
		serverConfig.DatabaseType = DefaultServerConfig.DatabaseType
	}
	if serverConfig.DatabaseConnStr == "" {
		serverConfig.DatabaseConnStr = DefaultServerConfig.DatabaseConnStr
	}
	if serverConfig.HTTPServerPort == 0 {
		serverConfig.HTTPServerPort = DefaultServerConfig.HTTPServerPort
	}

	// convert to absolute path
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	serverConfig.FilePath = absFilePath

	// return config
	return &serverConfig, nil
}
