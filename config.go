package main

import (
	"os"

	defaults "github.com/creasty/defaults"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	// enables debug logging
	Debug bool `yaml:"debug"`

	QueryUrl struct {
		// URL to query to obtain current public IPv4 address
		IPv4 string `yaml:"v4" validate:"required,url" default:"https://api4.my-ip.io/ip.txt"`

		// URL to query to obtain current public IPv6 address
		IPv6 string `yaml:"v6" validate:"required,url" default:"https://api6.my-ip.io/ip.txt"`
	} `yaml:"queryUrls"`

	Cloudflare struct {
		ApiToken string `yaml:"apiToken" validate:"required"`
		Zone     string `yaml:"zone" validate:"required"`
	} `yaml:"cloudflare"`

	Hostname string `yaml:"hostname" validate:"required"`
}

func makeConfig() *Configuration {
	c := new(Configuration)
	c.loadDefaults()
	return c
}

func (config *Configuration) loadDefaults() {
	defaults.Set(config)
}

func (config *Configuration) loadFromFile(path string) error {
	log.Debug("Loading configuration from ", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	return nil
}
