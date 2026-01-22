// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"strings"

	"github.com/sirupsen/logrus"
)

/*
AttachmentHeader provides information that describes an attachment. It has information
such as the type of content, file name, etc...
*/
type AttachmentHeader struct {
	ContentType             string `json:"contentType"`
	MIMEVersion             string `json:"mimeVersion"`
	ContentTransferEncoding string `json:"contentTransferEncoding"`
	ContentDisposition      string `json:"contentDisposition"`
	FileName                string `json:"fileName"`
	Body                    string `json:"body"`

	Logger *logrus.Entry
}

/*
NewAttachmentHeader creates a new AttachmentHeader object
*/
func NewAttachmentHeader(contentType, mimeVersion, contentTransferEncoding, contentDisposition, fileName, body string, logger *logrus.Entry) *AttachmentHeader {
	return &AttachmentHeader{
		ContentType:             contentType,
		MIMEVersion:             mimeVersion,
		ContentTransferEncoding: contentTransferEncoding,
		ContentDisposition:      contentDisposition,
		FileName:                fileName,
		Body:                    body,
		Logger:                  logger,
	}
}

/*
Parse processes a set of attachment headers. Splits lines up and figures out what
header data goes into what structure key. Most headers follow this format:

Header-Name: Some value here\r\n
*/
func (attachmentHeader *AttachmentHeader) Parse(contents string) {
	var key string

	attachmentHeader.FileName = ""
	attachmentHeader.ContentType = ""
	attachmentHeader.ContentDisposition = ""
	attachmentHeader.ContentTransferEncoding = ""
	attachmentHeader.MIMEVersion = ""
	attachmentHeader.Body = ""

	headerBodySplit := strings.Split(contents, "\r\n\r\n")

	if len(headerBodySplit) < 2 {
		attachmentHeader.Logger.Debugf("Attachment has no body content")
	} else {
		attachmentHeader.Body = strings.Join(headerBodySplit[1:], "\r\n\r\n")
	}

	contents = headerBodySplit[0]

	/*
	 * Unfold and split the header into lines. Loop over each line
	 * and figure out what headers are present. Store them.
	 * Sadly some headers require special processing.
	 */
	contents = UnfoldHeaders(contents)
	splitHeader := strings.Split(contents, "\r\n")
	numLines := len(splitHeader)

	for index := 0; index < numLines; index++ {
		splitItem := strings.Split(splitHeader[index], ":")
		key = splitItem[0]

		switch strings.ToLower(key) {
		case "content-disposition":
			contentDisposition := strings.TrimSpace(strings.Join(splitItem[1:], ""))
			attachmentHeader.Logger.WithField("contentDisposition", contentDisposition).Debugf("Attachment Content-Disposition")

			contentDispositionSplit := strings.Split(contentDisposition, ";")
			contentDispositionRightSide := strings.TrimSpace(strings.Join(contentDispositionSplit[1:], ";"))

			if len(contentDispositionSplit) < 2 || (len(contentDispositionSplit) > 1 && len(strings.TrimSpace(contentDispositionRightSide)) <= 0) {
				attachmentHeader.ContentDisposition = contentDisposition
			} else {
				attachmentHeader.ContentDisposition = strings.TrimSpace(contentDispositionSplit[0])

				/*
				 * See if we have an attachment and filename
				 */
				if strings.Contains(strings.ToLower(attachmentHeader.ContentDisposition), "attachment") && len(strings.TrimSpace(contentDispositionRightSide)) > 0 {
					filenameSplit := strings.Split(contentDispositionRightSide, "=")
					attachmentHeader.FileName = strings.Replace(strings.Join(filenameSplit[1:], "="), "\"", "", -1)
				}
			}

		case "content-transfer-encoding":
			attachmentHeader.ContentTransferEncoding = strings.TrimSpace(strings.Join(splitItem[1:], ""))
			attachmentHeader.Logger.WithField("content-transfer-encoding", attachmentHeader.ContentTransferEncoding).Debugf("Attachment Content-Transfer-Encoding")

		case "content-type":
			contentType := strings.TrimSpace(strings.Join(splitItem[1:], ""))
			attachmentHeader.Logger.WithField("content-type", contentType).Debugf("Attachment Content-Type")

			contentTypeSplit := strings.Split(contentType, ";")

			if len(contentTypeSplit) < 2 {
				attachmentHeader.ContentType = contentType
			} else {
				attachmentHeader.ContentType = strings.TrimSpace(contentTypeSplit[0])
				contentTypeRightSide := strings.TrimSpace(strings.Join(contentTypeSplit[1:], ";"))

				/*
				 * See if there is a "name" portion to this
				 */
				if strings.Contains(strings.ToLower(contentTypeRightSide), "name") || strings.Contains(strings.ToLower(contentTypeRightSide), "filename") {
					filenameSplit := strings.Split(contentTypeRightSide, "=")
					attachmentHeader.FileName = strings.Replace(strings.Join(filenameSplit[1:], "="), "\"", "", -1)
				}
			}

		case "mime-version":
			attachmentHeader.MIMEVersion = strings.TrimSpace(strings.Join(splitItem[1:], ""))
			attachmentHeader.Logger.WithField("mime-version", attachmentHeader.MIMEVersion).Debugf("Attachment MIME-Version")
		}
	}
}
