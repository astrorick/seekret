package seekret

type Config struct {
	configPath string
	serverPort uint16
}

func loadServerConfig(configPath string) (*Config, error) {
	// TODO: load config file 'filePath' and/or assign defaults

	// placeholder
	serverConfig := &Config{
		configPath: configPath,
		serverPort: 3553,
	}

	return serverConfig, nil
}
