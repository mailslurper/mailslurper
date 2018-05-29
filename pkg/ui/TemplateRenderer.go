// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo"

	"github.com/mailslurper/mailslurper/cmd/mailslurper/www"
)

var templates map[string]*template.Template

/*
TemplateRenderer describes a handlers for rendering layouts/pages
*/
type TemplateRenderer struct {
	templates *template.Template
}

/*
NewTemplateRenderer creates a new struct
*/
func NewTemplateRenderer(debugMode bool) *TemplateRenderer {
	result := &TemplateRenderer{}
	result.LoadTemplates(debugMode)

	return result
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	var tmpl *template.Template
	var ok bool

	if tmpl, ok = templates[name]; !ok {
		return fmt.Errorf("Cannot find template %s", name)
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

func (t *TemplateRenderer) LoadTemplates(debugMode bool) {
	templates = make(map[string]*template.Template)

	templates["mainLayout:admin"], _ = template.Must(
		template.New("layout").Parse(www.FSMustString(debugMode, "/www/mailslurper/layouts/mainLayout.gohtml")),
	).Parse(www.FSMustString(debugMode, "/www/mailslurper/pages/admin.gohtml"))

	templates["mainLayout:index"], _ = template.Must(
		template.New("layout").Parse(www.FSMustString(debugMode, "/www/mailslurper/layouts/mainLayout.gohtml")),
	).Parse(www.FSMustString(debugMode, "/www/mailslurper/pages/index.gohtml"))

	templates["mainLayout:manageSavedSearches"], _ = template.Must(
		template.New("layout").Parse(www.FSMustString(debugMode, "/www/mailslurper/layouts/mainLayout.gohtml")),
	).Parse(www.FSMustString(debugMode, "/www/mailslurper/pages/manageSavedSearches.gohtml"))

	templates["loginLayout:login"], _ = template.Must(
		template.New("layout").Parse(www.FSMustString(debugMode, "/www/mailslurper/layouts/loginLayout.gohtml")),
	).Parse(www.FSMustString(debugMode, "/www/mailslurper/pages/login.gohtml"))
}
