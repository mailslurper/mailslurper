package main

import "github.com/mailslurper/mailslurper/pkg/mailslurper"

func setupConfig() {
	/*
	 * Load configuration
	 */
	if config, err = mailslurper.LoadConfigurationFromFile(CONFIGURATION_FILE_NAME); err != nil {
		logger.WithError(err).Fatalf("There was an error reading the configuration file '%s'", CONFIGURATION_FILE_NAME)
	}
}
