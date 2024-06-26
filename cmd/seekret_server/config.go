package main

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// database parameters
	DatabaseType    string `yaml:"databaseType"`
	DatabaseConnStr string `yaml:"databaseConnStr"`

	// http server parameters
	HTTPServerPort uint16 `yaml:"httpServerPort"`

	// srp parameters
	SRPSaltSize uint64 `yaml:"srpSaltSize"`
	SRPHashFcn  string `yaml:"srpHashFcn"`

	// jwt parameters
	JWTSigningFcn string `yaml:"jwtSigningFcn"`
	JWTSigningKey string `yaml:"jwtSigningKey"`
}

// NewServerConfig tries to read and parse the provided file in a ServerConfig object.
// An error is returned if the provided file does not exist or if the file is not accessible.
func NewServerConfig(filePath string) (*Config, error) {
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
	var serverConfig Config
	if err := yaml.Unmarshal(configBytes, &serverConfig); err != nil {
		return nil, err
	}

	// return config
	return &serverConfig, nil
}
