package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/spf13/viper"
)

type Handlers []Handler
type Handler struct {
	Name        string
	ClusterName string
	WebhookURL  string
}

type Resources []Resource
type Resource struct {
	Kind     string
	Policies []string
}

type Config struct {
	Handlers  Handlers  `mapstructure:"handler"`
	Resources Resources `mapstructure:"resource"`
	CLI       *cli.CLIArgs
}

func readConfigFile(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func NewConfig() (*Config, error) {
	var err error
	c := cli.NewCLI()

	if err = readConfigFile(c.ConfigFile); err != nil {
		return nil, err
	}

	config := &Config{CLI: c}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	log.Debugf("Running kwatchman with config %#v", config)
	if config.Resources == nil || config.Handlers == nil {
		return nil, fmt.Errorf("Malformed %s config file", c.ConfigFile)
	}
	return config, nil
}
