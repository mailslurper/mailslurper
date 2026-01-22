// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"math"
	"regexp"
	"strings"
)

/*
An Attachment is any content embedded in the mail data that is not
considered the body
*/
type Attachment struct {
	ID       string            `json:"id"`
	MailID   string            `json:"mailId"`
	Headers  *AttachmentHeader `json:"headers"`
	Contents string            `json:"contents"`
}

/*
NewAttachment creates a new Attachment object
*/
func NewAttachment(id, mailid string, headers *AttachmentHeader, contents string) *Attachment {
	return &Attachment{
		ID:       id,
		MailID:   mailid,
		Headers:  headers,
		Contents: contents,
	}
}

/*
IsContentBase64 returns true/false if the content of this attachment
resembles a base64 encoded string.
*/
func (attachment *Attachment) IsContentBase64() bool {
	spaceKiller := func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}

		return r
	}

	trimmedContents := strings.Map(spaceKiller, attachment.Contents)

	if math.Mod(float64(len(trimmedContents)), 4.0) == 0 {
		matchResult, err := regexp.Match("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$", []byte(trimmedContents))
		if err == nil {
			if matchResult {
				return true
			}
		}
	}

	return false
}
