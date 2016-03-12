package layout

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mailslurper/mailslurper/global"
	"github.com/mailslurper/mailslurper/model"
	"github.com/mailslurper/mailslurper/www"
)

func RenderMainLayout(writer http.ResponseWriter, request *http.Request, htmlFileName string, data model.Page) error {
	var layout string
	var err error
	var tmpl *template.Template
	var pageString string

	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")

	/*
	 * Pre-load layout information
	 */
	if global.DEBUG_ASSETS {
		var bytes []byte

		if bytes, err = ioutil.ReadFile("./www/mailslurper/layouts/mainLayout.html"); err != nil {
			log.Printf("MailSlurper: ERROR - Error setting up layout: %s\n", err.Error())
			os.Exit(1)
		}

		layout = string(bytes)
	} else {
		if layout, err = www.FSString(false, "/www/mailslurper/layouts/mainLayout.html"); err != nil {
			log.Printf("MailSlurper: ERROR - Error setting up layout: %s\n", err.Error())
			os.Exit(1)
		}
	}

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
