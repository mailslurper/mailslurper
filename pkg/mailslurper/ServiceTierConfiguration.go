// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package mailslurper

import (
	"github.com/adampresley/webframework/logging"
)

/*
ServiceTierConfiguration allows a caller to configure how to start
and run the service tier HTTP server
*/
type ServiceTierConfiguration struct {
	Address          string
	Database         IStorage
	Log              *logging.Logger
	Port             int
	CertFile         string
	KeyFile          string
	CertIsSelfSigned bool
}
