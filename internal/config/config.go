package config

import (
	"fmt"
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

// String provided a formatted representation of the reference 'ServerConfig' object.
func (sc *ServerConfig) String() string {
	if sc.FilePath == "" {
		return fmt.Sprintf("Default Config: {\n\tDatabase Type: %s\n\tDatabase Connection String: \"%s\"\n\tHTTP Server Port: %d\n}", sc.DatabaseType, sc.DatabaseConnStr, sc.HTTPServerPort)
	} else {
		return fmt.Sprintf("From '%s': {\n\tDatabase Type: %s\n\tDatabase Connection String: \"%s\"\n\tHTTP Server Port: %d\n}", sc.FilePath, sc.DatabaseType, sc.DatabaseConnStr, sc.HTTPServerPort)
	}
}

// NewServerConfig tries to open, read and parse the provided file in a 'ServerConfig' object.
// If an empty 'filePath' is provided, it returns the default config.
func NewServerConfig(filePath string) (*ServerConfig, error) {
	// this will be returned when no 'filePath' is provided
	defaultConfig := &ServerConfig{
		FilePath:        "",
		DatabaseType:    "sqlite3",
		DatabaseConnStr: "seekret.db",
		HTTPServerPort:  3553,
	}

	// return default config if no path is provided
	if filePath == "" {
		return defaultConfig, nil
	}

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
		serverConfig.DatabaseType = defaultConfig.DatabaseType
	}
	if serverConfig.DatabaseConnStr == "" {
		serverConfig.DatabaseConnStr = defaultConfig.DatabaseConnStr
	}
	if serverConfig.HTTPServerPort == 0 {
		serverConfig.HTTPServerPort = defaultConfig.HTTPServerPort
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
