package controllers

import (
	"fmt"
	"net/http"

	"github.com/mailslurper/mailslurper/global"
)

/*
GetVersion outputs the current running version of this MailSlurper server instance
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, global.SERVER_VERSION)
}
