// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/mailslurper/libmailslurper/server"
	"github.com/mailslurper/mailslurper/global"
)

/*
GetVersion outputs the current running version of this MailSlurper server instance
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	result := server.Version{
		Version: global.SERVER_VERSION,
	}

	GoHttpService.WriteJson(writer, result, 200)
}
