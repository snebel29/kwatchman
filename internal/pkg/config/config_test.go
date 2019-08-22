package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"reflect"
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
	if len(config.Handlers) != 4 {
		t.Errorf("config.Handlers should have 4 item and has %d instead", len(config.Handlers))
	}
	found := false
	for _, h := range config.Handlers {
		if h.Name == "ignoreEvents" &&
			reflect.DeepEqual(h.IgnoreEvents, []string{"Add", "Delete"}) {
			found = true
		}
	}
	if !found {
		t.Errorf("Events should have been found, got %#v instead", config.Handlers)
	}

}

func TestNonExistantConfigShouldReturnError_NewConfig(t *testing.T) {
	fixture := "nonexistent.toml"
	_, err := loadConfigFileHelper(fixture)
	if err == nil {
		t.Error("there should have been an error")
	}
}

func TestWrongLogLevelShouldFail(t *testing.T) {
	os.Args = []string{
		"kwatchman",
		"--log-level=wrongLogLevel",
	}
	if _, err := NewConfig(); err == nil {
		t.Error("an error should have been returned")
	}
}

func TestGoodLogLevelShouldSet(t *testing.T) {
	level := "debug"
	configFile := path.Join(path.Dir(thisFilename), "fixtures", "config.toml")

	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--log-level=%s", level),
		fmt.Sprintf("--config=%s", configFile),
	}
	if _, err := NewConfig(); err != nil {
		t.Error(err)
	}

	if logrus.GetLevel().String() != level {
		t.Errorf("level %s should be set, got %s instead", level, logrus.GetLevel())
	}
}
