// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package homeController

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adampresley/GoHttpService"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadFile("./www/index.html")
	if err != nil {
		GoHttpService.Error(writer, "There was an error loading the index page")
		return
	}

	fmt.Fprintf(writer, string(body))
}
