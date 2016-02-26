// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/gorilla/context"
	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/mailslurper/model"
)

/*
GetServiceSettings returns the settings necessary to talk to the MailSlurper
back-end service tier.
*/
func GetServiceSettings(writer http.ResponseWriter, request *http.Request) {
	config := (context.Get(request, "config")).(*configuration.Configuration)

	settings := model.ServiceSettings{
		ServiceAddress: config.ServiceAddress,
		ServicePort:    config.ServicePort,
		Version:        "v1",
	}

	GoHttpService.WriteJson(writer, settings, 200)
}
