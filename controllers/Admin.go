// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"log"
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/gorilla/context"
)

/*
Admin is the page for performing administrative tasks in MailSlurper
*/
func Admin(writer http.ResponseWriter, request *http.Request) {
	layout := (context.Get(request, "layout")).(GoHttpService.Layout)

	data := struct {
		Title string
	}{
		"Admin",
	}

	err := layout.RenderView(writer, "admin", data)
	if err != nil {
		log.Println("MailSlurper: ERROR - Problem rendering view 'admin' -", err.Error())
		GoHttpService.Error(writer, "There was an error retrieving and rendering the page 'admin'")
		return
	}
}
