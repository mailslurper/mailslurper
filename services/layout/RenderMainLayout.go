package layout

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/mailslurper/mailslurper/global"
	"github.com/mailslurper/mailslurper/www"

	"github.com/gorilla/context"
)

func RenderMainLayout(writer http.ResponseWriter, request *http.Request, htmlFileName string, data interface{}) error {
	layout := (context.Get(request, "layout")).(string)
	var err error
	var tmpl *template.Template
	var pageString string

	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if tmpl, err = template.New("layout").Parse(layout); err != nil {
		return err
	}

	if pageString, err = getHTMLPageString(htmlFileName); err != nil {
		return err
	}

	if tmpl, err = tmpl.Parse(pageString); err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func getHTMLPageString(htmlFileName string) (string, error) {
	if global.DEBUG_ASSETS {
		bytes, err := ioutil.ReadFile("./www/" + htmlFileName)
		return string(bytes), err
	}

	return www.FSString(false, "/www/"+htmlFileName)
}
