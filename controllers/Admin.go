// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/mailslurper/mailslurper/model"
	"github.com/mailslurper/mailslurper/services/layout"
)

/*
Admin is the page for performing administrative tasks in MailSlurper
*/
func Admin(writer http.ResponseWriter, request *http.Request) {
	var err error

	data := model.Page{
		Title: "Admin",
	}

	if err = layout.RenderMainLayout(writer, request, "admin.html", data); err != nil {
		GoHttpService.Error(writer, err.Error())
	}
}
