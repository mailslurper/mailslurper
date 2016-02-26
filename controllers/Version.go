// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"fmt"
	"net/http"

	"github.com/mailslurper/mailslurper/global"
)

/*
GetVersion outputs the current running version of this MailSlurper server instance
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, global.SERVER_VERSION)
}
