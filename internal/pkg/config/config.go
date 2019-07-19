package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/spf13/viper"
)

// Handlers holds a list of Handler
type Handlers []Handler

// Handler struct holds configuration fields for all handlers in this code base,
// non configured options will be set to its zero value in the struct filed, is the responsability
// of the handler to validate them, as well as to avoid naming conflicts, in the future namespacing
// might me implemented
type Handler struct {
	Name         string   // Used by all handlers
	ClusterName  string   // Used by slack handler
	WebhookURL   string   // Used by slack handler
	IgnoreEvents []string // Used by ignoreEvents handler
}

// Resources holds a list of Resource
type Resources []Resource

// Resource holds the individual resource configurations
type Resource struct {
	Kind string
}

// Config represent the config file
type Config struct {
	Handlers  Handlers  `mapstructure:"handler"`
	Resources Resources `mapstructure:"resource"`
	CLI       *cli.Args
}

func readConfigFile(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// NewConfig return kwatchman config file unmarshalled using viper
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
		return nil, fmt.Errorf("malformed %s config file", c.ConfigFile)
	}
	return config, nil
}
