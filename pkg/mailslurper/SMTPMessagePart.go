// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"bufio"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
)

/*
An SMTPMessagePart represents a single message/content from a DATA transmission
from an SMTP client. This contains the headers and body content. It also contains
a reference to a collection of sub-messages, if any. This allows us to support
the recursive tree-like nature of the MIME protocol.
*/
type SMTPMessagePart struct {
	Message      *mail.Message
	MessageParts []ISMTPMessagePart

	body   string
	logger *logrus.Entry
}

/*
NewSMTPMessagePart returns a new instance of this struct
*/
func NewSMTPMessagePart(logger *logrus.Entry) *SMTPMessagePart {
	return &SMTPMessagePart{
		Message:      &mail.Message{},
		MessageParts: make([]ISMTPMessagePart, 0),

		body:   "",
		logger: logger,
	}
}

/*
AddBody adds body content
*/
func (messagePart *SMTPMessagePart) AddBody(body string) error {
	messagePart.Message.Body = strings.NewReader(body)
	return nil
}

/*
AddHeaders takes a header set and adds it to this message part.
*/
func (messagePart *SMTPMessagePart) AddHeaders(headers textproto.MIMEHeader) error {
	messagePart.Message.Header = mail.Header(headers)
	return nil
}

/*
BuildMessages pulls the message body from the data transmission
and stores the whole body. If the message type is multipart it then
attempts to parse the parts.
*/
func (messagePart *SMTPMessagePart) BuildMessages(body string) error {
	var err error
	var isMultipart bool
	var boundary string
	var headers textproto.MIMEHeader

	headerReader := textproto.NewReader(bufio.NewReader(strings.NewReader(body)))

	if headers, err = headerReader.ReadMIMEHeader(); err != nil {
		return errors.Wrap(err, "Problem reading headers")
	}

	messagePart.AddHeaders(headers)

	/*
	 * If this is not a multipart message, bail early. We've got
	 * what we need.
	 */
	if isMultipart, err = messagePart.ContentIsMultipart(); err != nil {
		return errors.Wrapf(err, "Error getting content type information in message part")
	}

	if !isMultipart {
		messagePart.logger.Debugf("Body of message: %s", body)

		if err = messagePart.AddBody(body); err != nil {
			return errors.Wrapf(err, "Error adding body to message part")
		}

		return nil
	}

	if boundary, err = messagePart.GetBoundary(); err != nil {
		return errors.Wrapf(err, "Error getting boundary for message part")
	}

	messagePart.logger.Debugf("Body of message: %s", body)
	if err = messagePart.AddBody(body); err != nil {
		return errors.Wrapf(err, "Error adding body to message part")
	}

	return messagePart.ParseMessages(body, boundary)
}

/*
GetBody retrieves the body portion of the message
*/
func (messagePart *SMTPMessagePart) GetBody() string {
	var err error
	var bytes []byte

	if messagePart.body == "" {
		if bytes, err = ioutil.ReadAll(messagePart.Message.Body); err != nil {
			messagePart.logger.Errorf("Problem reading message body: %s", err.Error())
			return ""
		}

		messagePart.body = string(bytes)
	}

	return messagePart.body
}

/*
GetFilenameFromContentDisposition returns a filename from a Content-Disposition header
*/
func (messagePart *SMTPMessagePart) GetFilenameFromContentDisposition() string {
	contentDisposition := messagePart.GetContentDisposition()
	contentDispositionSplit := strings.Split(contentDisposition, ";")
	contentDispositionRightSide := strings.TrimSpace(strings.Join(contentDispositionSplit[1:], ";"))

	fileName := ""

	if strings.Contains(strings.ToLower(contentDisposition), "attachment") && len(strings.TrimSpace(contentDispositionRightSide)) > 0 {
		filenameSplit := strings.Split(contentDispositionRightSide, "=")
		fileName = strings.Replace(strings.Join(filenameSplit[1:], "="), "\"", "", -1)
	}

	return fileName
}

/*
GetHeader returns the value of a specified header key
*/
func (messagePart *SMTPMessagePart) GetHeader(key string) string {
	decoder := new(mime.WordDecoder)
	decoder.CharsetReader = func(headerCharset string, input io.Reader) (io.Reader, error) {
		encoding, _ := charset.Lookup(headerCharset)
		return encoding.NewDecoder().Reader(input), nil
	}

	result, _ := decoder.DecodeHeader(messagePart.Message.Header.Get(key))
	return result
}

/*
GetMessageParts returns any additional sub-messages related to this message
*/
func (messagePart *SMTPMessagePart) GetMessageParts() []ISMTPMessagePart {
	return messagePart.MessageParts
}

/*
ParseMessages parses messages in an SMTP body
*/
func (messagePart *SMTPMessagePart) ParseMessages(body string, boundary string) error {
	var err error
	var bodyPart []byte
	var part *multipart.Part

	reader := multipart.NewReader(strings.NewReader(body), boundary)

	for {
		part, err = reader.NextPart()

		switch err {
		case io.EOF:
			messagePart.logger.Debugf("Reached EOF for part")
			return nil

		case nil:
			if bodyPart, err = ioutil.ReadAll(part); err != nil {
				return errors.Wrapf(err, "Error reading body for content type '%s'", messagePart.Message.Header.Get("Content-Type"))
			}

			innerBody := string(bodyPart)
			messagePart.logger.Debugf("Building new message part: %v", part.Header)

			if boundary, err = messagePart.GetBoundaryFromHeaderString(part.Header.Get("Content-Type")); err != nil {
				return errors.Wrapf(err, "Error getting boundary marker")
			}

			newMessage := NewSMTPMessagePart(messagePart.logger)
			newMessage.Message.Header = messagePart.convertPartHeadersToMap(part.Header)
			newMessage.Message.Body = strings.NewReader(innerBody)

			newMessage.ParseMessages(innerBody, boundary)
			messagePart.MessageParts = append(messagePart.MessageParts, newMessage)

		default:
			return errors.Wrapf(err, "Error reading next part for content type '%s'", messagePart.Message.Header.Get("Content-Type"))
		}
	}
}

/*
ContentIsMultipart returns true if the Content-Type header contains "multipart"
*/
func (messagePart *SMTPMessagePart) ContentIsMultipart() (bool, error) {
	mediaType, _, err := messagePart.parseContentType()
	return strings.HasPrefix(mediaType, "multipart/"), err
}

/*
GetBoundary returns the message boundary string
*/
func (messagePart *SMTPMessagePart) GetBoundary() (string, error) {
	_, boundary, err := messagePart.parseContentType()
	return boundary, err
}

/*
GetBoundaryFromHeaderString returns the boundary marker defined in the header
*/
func (messagePart *SMTPMessagePart) GetBoundaryFromHeaderString(header string) (string, error) {
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return "", err
	}

	return params["boundary"], nil
}

/*
GetContentDisposition returns the value of the Content-Disposition header
*/
func (messagePart *SMTPMessagePart) GetContentDisposition() string {
	return messagePart.Message.Header.Get("Content-Disposition")
}

/*
GetContentType returns the value from the Content-Type header
*/
func (messagePart *SMTPMessagePart) GetContentType() string {
	return messagePart.Message.Header.Get("Content-Type")
}

func (messagePart *SMTPMessagePart) parseContentType() (string, string, error) {
	contentType := messagePart.GetContentType()
	if contentType == "" {
		return "", "", nil
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", "", err
	}

	return mediaType, params["boundary"], nil
}

func (messagePart *SMTPMessagePart) convertPartHeadersToMap(partHeaders textproto.MIMEHeader) map[string][]string {
	convertedHeaders := make(map[string][]string)
	for key, value := range partHeaders {
		convertedHeaders[key] = value
	}

	return convertedHeaders
}
