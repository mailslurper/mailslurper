// Copyright 2013-3014 Adam Presley. All rights reserved
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
ManageSavedSearches is the page for managing saved searches
*/
func ManageSavedSearches(writer http.ResponseWriter, request *http.Request) {
	layout := (context.Get(request, "layout")).(GoHttpService.Layout)

	data := struct {
		Title string
	}{
		"Manage Saved Searches",
	}

	err := layout.RenderView(writer, "manageSavedSearches", data)
	if err != nil {
		log.Println("MailSlurper: ERROR - Problem rendering view 'manageSavedSearches' -", err.Error())
		GoHttpService.Error(writer, "There was an error retrieving and rendering the page 'manageSavedSearches'")
		return
	}
}
