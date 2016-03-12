// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/mailslurper/mailslurper/services/layout"
)

/*
Index is the main view. This endpoint provides the email list and email detail
views.
*/
func Index(writer http.ResponseWriter, request *http.Request) {
	var err error

	data := struct {
		Title string
	}{
		"Mail",
	}

	if err = layout.RenderMainLayout(writer, request, "index.html", data); err != nil {
		GoHttpService.Error(writer, err.Error())
	}
}
