// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package settingsController

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/mailslurper/mailslurper/global"
)

func GetServiceSettings(writer http.ResponseWriter, request *http.Request) {
	data := struct {
		ServiceAddress string `json:"serviceAddress"`
		ServicePort    int    `json:"servicePort"`
	}{
		global.Config.ServiceAddress,
		global.Config.ServicePort,
	}

	GoHttpService.WriteJson(writer, data, 200)
}
