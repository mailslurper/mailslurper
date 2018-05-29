// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailslurper/mailslurper/pkg/mailslurper"
	"github.com/skratchdot/open-golang/open"
)

/*
StartBrowser opens the user's default browser to the configured URL
*/
func StartBrowser(config *mailslurper.Configuration, logger *logrus.Entry) {
	timer := time.NewTimer(time.Second)

	go func() {
		<-timer.C
		logger.Infof("Opening web browser to http://%s:%d", config.WWWAddress, config.WWWPort)
		err := open.Start(fmt.Sprintf("http://%s:%d", config.WWWAddress, config.WWWPort))
		if err != nil {
			logger.Infof("ERROR - Could not open browser - %s", err.Error())
		}
	}()
}
