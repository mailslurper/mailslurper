// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
ServiceSettings represents the necessary settings to connect to
and talk to the MailSlurper service tier.
*/
type ServiceSettings struct {
	ServiceAddress string `json:"serviceAddress"`
	ServicePort    int    `json:"servicePort"`
	Version        string `json:"version"`
}
