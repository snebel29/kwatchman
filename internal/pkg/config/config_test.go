package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

var thisFilename string

func init() {
	_, t, _, _ := runtime.Caller(0)
	thisFilename = t
}

func TestReadConfigFile(t *testing.T) {
	err := readConfigFile("nonExistentFile")
	if err == nil {
		t.Error("Non existent file should have returned an error")
	}
}

func loadConfigFileHelper(fixture string) (*Config, error) {
	configFile := path.Join(path.Dir(thisFilename), "fixtures", fixture)
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--config=%s", configFile),
	}
	return NewConfig()
}

func malformedConfigFileHelper(fixture string, t *testing.T) {
	_, err := loadConfigFileHelper(fixture)
	if err == nil {
		t.Errorf("%s file should have returned an error", fixture)
	}
}

func TestNewConfigReturnErrorWhenFileisMalformed(t *testing.T) {
	malformedConfigFileHelper("handlerless-config.toml", t)
	malformedConfigFileHelper("resourcesless-config.toml", t)
}

func TestGoodNewConfigShouldParseCorrectly(t *testing.T) {
	fixture := "config.toml"
	config, err := loadConfigFileHelper(fixture)
	if err != nil {
		t.Errorf("%s file should have NOT returned an error: %s", fixture, err)
	}
	if len(config.Resources) != 1 {
		t.Errorf("config.Resources should have 1 item and has %d instead", len(config.Resources))
	}
	if len(config.Handlers) != 3 {
		t.Errorf("config.Resources should have 1 item and has %d instead", len(config.Handlers))
	}
}
