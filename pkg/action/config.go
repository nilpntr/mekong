package action

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Config struct {
	// The host of the backend to be redirected to
	ListenPorts []string           `yaml:"listenPorts"`
	Routes      []ConfigRoutes     `yaml:"routes"`
	Heartbeat   *[]ConfigHeartbeat `yaml:"heartbeat,omitempty"`
}

type ConfigHeartbeat struct {
	Path         string  `yaml:"path"`
	ResponseCode *int    `yaml:"response_code,omitempty"`
	Message      *string `yaml:"message,omitempty"`
}

type ConfigRoutes struct {
	Path                string                `yaml:"path"`
	BackendHost         string                `yaml:"backendHost"`
	Methods             []HTTPMethod          `yaml:"methods"`
	BasicAuthentication *ConfigRouteBasicAuth `yaml:"basicAuthentication,omitempty"`
	Rules               *ConfigRouteRules     `yaml:"rules,omitempty"`
	Headers             *[]string             `yaml:"headers,omitempty"`
	Debug               *ConfigRouteDebug     `yaml:"debug,omitempty"`
}

type ConfigRouteDebug struct {
	Headers *bool `yaml:"headers,omitempty"`
	Body    *bool `yaml:"body,omitempty"`
	URL     *bool `yaml:"url,omitempty"`
}

type ConfigRouteBasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ConfigRouteRules struct {
	HasQueryString *bool `yaml:"hasQueryString,omitempty"`
	HasBody        *bool `yaml:"hasBody,omitempty"`
}

func newConfig(configPath string) (*Config, error) {
	config := new(Config)

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(ErrConfigFileNotFound)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		log.Fatal(err)
	}

	return config, nil
}

func validateConfig(config *Config) error {
	for _, elem := range config.Routes {
		if elem.BackendHost == "" {
			return ErrValidateFieldNotExistsBackendHost
		}
	}

	for _, elem := range config.ListenPorts {
		if !strings.HasPrefix(elem, ":") {
			return fmt.Errorf("the following listen port %s needs to start with a colon", elem)
		}
	}

	return nil
}
