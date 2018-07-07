// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import "github.com/mailslurper/mailslurper/pkg/mailslurper"
import "os"

func setupConfig() {
	var err error
	var configFile string

	if len(os.Getenv("MS_CONFIG_DIR")) > 0 {
		configFile = os.Getenv("MS_CONFIG_DIR") + "/" + CONFIGURATION_FILE_NAME
	} else {
		configFile = CONFIGURATION_FILE_NAME
	}

	/*
	 * Load configuration
	 */
	if config, err = mailslurper.LoadConfigurationFromFile(configFile); err != nil {
		logger.WithError(err).Fatalf("There was an error reading the configuration file '%s'", configFile)
	}
}
