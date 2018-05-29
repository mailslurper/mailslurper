// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import "github.com/mailslurper/mailslurper/pkg/mailslurper"

func setupConfig() {
	var err error

	/*
	 * Load configuration
	 */
	if config, err = mailslurper.LoadConfigurationFromFile(CONFIGURATION_FILE_NAME); err != nil {
		logger.WithError(err).Fatalf("There was an error reading the configuration file '%s'", CONFIGURATION_FILE_NAME)
	}
}
