// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/gorilla/context"
)

/*
Index is the main view. Here you look at mail.
*/
func Index(writer http.ResponseWriter, request *http.Request) {
	layout := (context.Get(request, "layout")).(GoHttpService.Layout)

	data := struct {
		Title string
	}{
		"Mail",
	}

	err := layout.RenderView(writer, "index", data)
	if err != nil {
		GoHttpService.Error(writer, err.Error())
		return
	}
}
