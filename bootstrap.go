package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type BootstrapInfo struct {
	ConfigFile string
	Debug      bool
}

func parseCommandLineArgs() BootstrapInfo {
	info := BootstrapInfo{}
	flag.StringVar(&info.ConfigFile, "file", "/etc/ip-updater.yaml", "The location of the configuration file to load")
	flag.BoolVar(&info.Debug, "debug", false, "Enables debug logging")
	flag.Parse()
	return info
}

func setLogLevel(debug bool) {
	// set log level
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func bootstrap() (*Configuration, error) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	// get bootstrap instructions from command line args
	bootstrapInfo := parseCommandLineArgs()
	setLogLevel(bootstrapInfo.Debug)

	// copy settings into config object
	config := makeConfig()
	config.Debug = bootstrapInfo.Debug

	// load config file
	err := config.loadFromFile(bootstrapInfo.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading config file %s: %w", bootstrapInfo.ConfigFile, err)
	}

	// re-apply settings if they changed
	if bootstrapInfo.Debug != config.Debug {
		setLogLevel(config.Debug)
	}

	if config.Debug {
		log.Debugf("Current configuration: %+v", config)
	}

	return config, nil
}
