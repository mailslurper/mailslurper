// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
Version is a struct that holds information about the most current
version of MailSlurper
*/
type Version struct {
	Version string `json:"version"`
}

/*
GetServerVersionFromMaster retieves the latest version information
for MailSlurper from the version.json file at master in the
origin repository
*/
func GetServerVersionFromMaster() (*Version, error) {
	var result *Version

	client := http.Client{}
	response, err := client.Get("https://raw.githubusercontent.com/mailslurper/mailslurper/master/cmd/mailslurper/version.json")

	if err != nil {
		return result, err
	}

	versionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(versionBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}
