package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/p0vidl0/mylslurper/internal/idgen"
)

// Parse turns a raw SMTP DATA payload plus its envelope addresses into a
// fully parsed Item, walking any MIME multipart structure to separate the
// text body, HTML body, and attachments.
func Parse(raw string, envelopeFrom string, envelopeTo []string) (*Item, error) {
	msg, err := mail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("parsing message headers: %w", err)
	}

	dec := new(mime.WordDecoder)
	decodeHeader := func(name string) string {
		v := msg.Header.Get(name)
		if v == "" {
			return v
		}
		if decoded, err := dec.DecodeHeader(v); err == nil {
			return decoded
		}
		return v
	}

	item := &Item{
		ID:         idgen.New(),
		From:       envelopeFrom,
		To:         envelopeTo,
		Subject:    decodeHeader("Subject"),
		XMailer:    msg.Header.Get("X-Mailer"),
		RawMessage: raw,
	}

	if item.From == "" {
		item.From = msg.Header.Get("From")
	}
	if len(item.To) == 0 {
		if to := msg.Header.Get("To"); to != "" {
			item.To = []string{to}
		}
	}

	if dateHeader := msg.Header.Get("Date"); dateHeader != "" {
		if t, err := mail.ParseDate(dateHeader); err == nil {
			item.DateSent = t
		}
	}
	if item.DateSent.IsZero() {
		item.DateSent = time.Now().UTC()
	}

	contentType := msg.Header.Get("Content-Type")
	item.ContentType = contentType

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// No/invalid Content-Type header: treat the whole body as plain text.
		body, _ := io.ReadAll(msg.Body)
		item.TextBody = string(body)
		return item, nil
	}
	item.Boundary = params["boundary"]

	if strings.HasPrefix(mediaType, "multipart/") {
		if err := walkMultipart(item, msg.Body, params["boundary"]); err != nil {
			return nil, err
		}
		return item, nil
	}

	body, err := decodeBody(msg.Body, textproto.MIMEHeader(msg.Header))
	if err != nil {
		return nil, err
	}
	if mediaType == "text/html" {
		item.HTMLBody = string(body)
	} else {
		item.TextBody = string(body)
	}

	return item, nil
}

// walkMultipart recursively walks a multipart body, filling in item's text
// body, HTML body, and attachments as it finds them.
func walkMultipart(item *Item, r io.Reader, boundary string) error {
	if boundary == "" {
		return fmt.Errorf("multipart message is missing a boundary")
	}

	mr := multipart.NewReader(r, boundary)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading multipart body: %w", err)
		}

		if err := processPart(item, part); err != nil {
			return err
		}
	}
}

func processPart(item *Item, part *multipart.Part) error {
	defer part.Close()

	header := textproto.MIMEHeader(part.Header)
	contentType := header.Get("Content-Type")

	mediaType := "text/plain"
	var params map[string]string
	if contentType != "" {
		var err error
		mediaType, params, err = mime.ParseMediaType(contentType)
		if err != nil {
			mediaType = "text/plain"
		}
	}

	fileName := attachmentFileName(header, params)

	if strings.HasPrefix(mediaType, "multipart/") {
		return walkMultipart(item, part, params["boundary"])
	}

	body, err := decodeBody(part, header)
	if err != nil {
		return err
	}

	if fileName != "" {
		item.Attachments = append(item.Attachments, Attachment{
			ID:          idgen.New(),
			MailItemID:  item.ID,
			FileName:    fileName,
			ContentType: contentType,
			Content:     body,
			Size:        len(body),
		})
		return nil
	}

	switch mediaType {
	case "text/html":
		item.HTMLBody += string(body)
	default:
		item.TextBody += string(body)
	}

	return nil
}

// attachmentFileName returns the file name for a part if it should be
// treated as an attachment, or "" if it's an inline text/html body part.
func attachmentFileName(header textproto.MIMEHeader, contentTypeParams map[string]string) string {
	if disposition := header.Get("Content-Disposition"); disposition != "" {
		if _, params, err := mime.ParseMediaType(disposition); err == nil {
			if name := decodeFileName(params["filename"]); name != "" {
				return name
			}
			if strings.EqualFold(strings.SplitN(disposition, ";", 2)[0], "attachment") {
				return "attachment"
			}
		}
	}

	if name := decodeFileName(contentTypeParams["name"]); name != "" {
		return name
	}

	return ""
}

func decodeFileName(raw string) string {
	if raw == "" {
		return ""
	}
	dec := new(mime.WordDecoder)
	if decoded, err := dec.DecodeHeader(raw); err == nil {
		return decoded
	}
	return raw
}

// decodeBody reads r fully and reverses any Content-Transfer-Encoding.
func decodeBody(r io.Reader, header textproto.MIMEHeader) ([]byte, error) {
	switch strings.ToLower(header.Get("Content-Transfer-Encoding")) {
	case "base64":
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		decoded, err := base64.StdEncoding.DecodeString(strings.Map(func(r rune) rune {
			if r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, string(data)))
		if err != nil {
			return nil, fmt.Errorf("decoding base64 body: %w", err)
		}
		return decoded, nil
	case "quoted-printable":
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, quotedprintable.NewReader(r)); err != nil {
			return nil, fmt.Errorf("decoding quoted-printable body: %w", err)
		}
		return buf.Bytes(), nil
	default:
		return io.ReadAll(r)
	}
}
