// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func getMailAndAttachmentsQuery(whereClause string) string {
	sqlQuery := `
		SELECT
			  mailitem.dateSent
			, mailitem.fromAddress
			, mailitem.toAddressList
			, mailitem.subject
			, mailitem.xmailer
			, mailitem.body
			, mailitem.contentType
			, mailitem.boundary
			, attachment.id AS attachmentID
			, attachment.fileName
			, attachment.contentType

		FROM mailitem
			LEFT JOIN attachment ON attachment.mailItemID=mailitem.id

		WHERE 1=1 `

	sqlQuery = sqlQuery + whereClause
	sqlQuery = sqlQuery + ` ORDER BY mailitem.dateSent DESC `

	return sqlQuery
}

func getMailCountQuery(mailSearch *MailSearch) (string, []interface{}) {
	sqlQuery := `
		SELECT COUNT(id) AS mailItemCount FROM mailitem WHERE 1=1
	`

	var parameters []interface{}
	return addSearchCriteria(sqlQuery, parameters, mailSearch)
}

func getDeleteAttachmentsQuery(startDate string) string {
	where := ""

	if len(startDate) > 0 {
		where = where + " AND mailitem.dateSent <= ? "
	}

	sqlQuery := "DELETE FROM attachment WHERE attachment.mailItemID IN (SELECT mailitem.id FROM mailitem WHERE 1=1 " + where + ")"
	return sqlQuery
}

func getDeleteMailQuery(startDate string) string {
	where := ""

	if len(startDate) > 0 {
		where = where + " AND dateSent <= ? "
	}

	sqlQuery := "DELETE FROM mailitem WHERE 1=1" + where
	return sqlQuery
}

func getInsertMailQuery() string {
	sqlQuery := `
		INSERT INTO mailitem (
			  id
			, dateSent
			, fromAddress
			, toAddressList
			, subject
			, xmailer
			, body
			, contentType
			, boundary
		) VALUES (
			  ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
			, ?
		)
	`

	return sqlQuery
}

func getInsertAttachmentQuery() string {
	sqlQuery := `
		INSERT INTO attachment (
			  id
			, mailItemId
			, fileName
			, contentType
			, content
		) VALUES (
			  ?
			, ?
			, ?
			, ?
			, ?
		)
	`

	return sqlQuery
}

func addOrderBy(sqlQuery string, tablePrefix string, mailSearch *MailSearch) string {
	switch mailSearch.OrderByField {
	case "subject":
		sqlQuery += fmt.Sprintf(" ORDER BY %s.subject ", tablePrefix)

	case "from":
		sqlQuery += fmt.Sprintf(" ORDER BY %s.fromAddress ", tablePrefix)

	default:
		sqlQuery += fmt.Sprintf(" ORDER BY %s.dateSent ", tablePrefix)
	}

	switch mailSearch.OrderByDirection {
	case "asc":
		sqlQuery += " ASC "

	default:
		sqlQuery += " DESC "
	}

	return sqlQuery
}

func addSearchCriteria(sqlQuery string, parameters []interface{}, mailSearch *MailSearch) (string, []interface{}) {
	var date time.Time
	var err error

	if len(strings.TrimSpace(mailSearch.Message)) > 0 {
		sqlQuery += `
			AND (
				mailitem.body LIKE ?
				OR mailitem.subject LIKE ?
			)
		`

		parameters = append(parameters, "%"+mailSearch.Message+"%")
		parameters = append(parameters, "%"+mailSearch.Message+"%")
	}

	if len(strings.TrimSpace(mailSearch.From)) > 0 {
		sqlQuery += `
			AND mailitem.fromAddress LIKE ?
		`

		parameters = append(parameters, "%"+mailSearch.From+"%")
	}

	if len(strings.TrimSpace(mailSearch.To)) > 0 {
		sqlQuery += `
			AND mailitem.toAddressList LIKE ?
		`

		parameters = append(parameters, "%"+mailSearch.To+"%")
	}

	if len(strings.TrimSpace(mailSearch.Start)) > 0 {
		if date, err = time.Parse("2006-01-02", mailSearch.Start); err == nil {
			sqlQuery += `
				AND mailitem.dateSent >= ?
			`

			parameters = append(parameters, date)
		}
	}

	if len(strings.TrimSpace(mailSearch.End)) > 0 {
		if date, err = time.Parse("2006-01-02", mailSearch.End); err == nil {
			datePlusOne := date.Add(time.Hour * 24)

			sqlQuery += `
				AND mailitem.dateSent < ?
			`

			parameters = append(parameters, datePlusOne)
		}
	}

	return sqlQuery, parameters
}

func storeAttachments(mailItemID string, transaction *sql.Tx, attachments []*Attachment) error {
	var err error
	var attachmentID string

	for _, currentAttachment := range attachments {
		if attachmentID, err = GenerateID(); err != nil {
			return fmt.Errorf("Error generating ID for attachment: %s", err.Error())
		}

		statement, err := transaction.Prepare(getInsertAttachmentQuery())
		if err != nil {
			return fmt.Errorf("Error preparing insert attachment statement: %s", err.Error())
		}

		_, err = statement.Exec(
			attachmentID,
			mailItemID,
			currentAttachment.Headers.FileName,
			currentAttachment.Headers.ContentType,
			currentAttachment.Contents,
		)

		if err != nil {
			return fmt.Errorf("Error executing insert attachment in StoreMail: %s", err.Error())
		}

		statement.Close()
		currentAttachment.ID = attachmentID
	}

	return nil
}
