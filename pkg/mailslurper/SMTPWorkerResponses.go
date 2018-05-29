// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
Responses that are sent to SMTP clients
*/
const (
	SMTP_CRLF                     string = "\r\n"
	SMTP_DATA_TERMINATOR          string = "\r\n.\r\n"
	SMTP_WELCOME_MESSAGE          string = "220 Welcome to MailSlurper!"
	SMTP_CLOSING_MESSAGE          string = "221 Bye"
	SMTP_OK_MESSAGE               string = "250 Ok"
	SMTP_DATA_RESPONSE_MESSAGE    string = "354 End data with <CR><LF>.<CR><LF>"
	SMTP_HELLO_RESPONSE_MESSAGE   string = "250 Hello. How very nice to meet you!"
	SMTP_ERROR_TRANSACTION_FAILED string = "554 Transaction failed"
)
