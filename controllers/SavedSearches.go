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
ManageSavedSearches is the page for managing saved searches
*/
func ManageSavedSearches(writer http.ResponseWriter, request *http.Request) {
	var err error

	data := struct {
		Title string
	}{
		"Manage Saved Searches",
	}

	if err = layout.RenderMainLayout(writer, request, "manageSavedSearches.html", data); err != nil {
		GoHttpService.Error(writer, err.Error())
	}
}
