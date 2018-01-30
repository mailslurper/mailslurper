package mailslurper_test

import (
	"fmt"
	"net/smtp"

	"github.com/mailslurper/mailslurper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mailslurper", func() {
	var address string
	var auth smtp.Auth

	BeforeEach(func() {
		address = "localhost:2500"
		auth = smtp.PlainAuth("", "", "", "adampresley.com")

		DeleteAllMail()
	})

	Describe("Sending a valid email", func() {
		var from string

		BeforeEach(func() {
			from = "adam@adampresley.com"
		})

		/*
		 * Valid plain text email
		 */
		Context("that is text/plain", func() {
			It("records the plain text in the database", func() {
				var err error
				var mailItems []mailslurper.MailItem

				body := "This is a plain text email"
				to := []string{"bob@test.com"}
				msg := []byte("To: bob@test.com\r\n" +
					"Subject: Plain Text Test\r\n" +
					"Date: Thu, 08 Dec 2016 23:46:05 -0600 CST\r\n" +
					"Content-Type: text/plain\r\n" +
					"\r\n" +
					body + "\r\n")

				smtp.SendMail(address, auth, from, to, msg)
				Pause()

				search := &mailslurper.MailSearch{}
				if mailItems, err = database.GetMailCollection(0, 1, search); err != nil {
					Fail(fmt.Sprintf("Error getting mail collection from database: %s", err.Error()))
				}

				Expect(len(mailItems)).To(Equal(1))
				Expect(mailItems[0].Subject).To(Equal("Plain Text Test"))
				Expect(mailItems[0].DateSent).To(Equal("2016-12-08 23:46:05"))
				Expect(mailItems[0].ContentType).To(Equal("text/plain"))
				Expect(mailItems[0].FromAddress).To(Equal(from))
				Expect(mailItems[0].ToAddresses).To(Equal([]string{
					"bob@test.com",
				}))
				Expect(mailItems[0].Body).To(Equal(body))
			})

			It("case accept 'data' in the 'TO' field", func() {
				var err error
				var mailItems []mailslurper.MailItem

				body := "This is an email with 'data' in the 'TO' field"
				to := []string{"data@test.com"}
				msg := []byte("To: data@test.com\r\n" +
					"Subject: 'Data' Test\r\n" +
					"Date: Thu, 08 Dec 2016 23:46:05 -0600 CST\r\n" +
					"Content-Type: text/plain\r\n" +
					"\r\n" +
					body + "\r\n")

				smtp.SendMail(address, auth, from, to, msg)
				Pause()

				search := &mailslurper.MailSearch{}
				if mailItems, err = database.GetMailCollection(0, 1, search); err != nil {
					Fail(fmt.Sprintf("Error getting mail collection from database: %s", err.Error()))
				}

				Expect(len(mailItems)).To(Equal(1))
				Expect(mailItems[0].Subject).To(Equal("&#39;Data&#39; Test"))
				Expect(mailItems[0].DateSent).To(Equal("2016-12-08 23:46:05"))
				Expect(mailItems[0].ContentType).To(Equal("text/plain"))
				Expect(mailItems[0].FromAddress).To(Equal(from))
				Expect(mailItems[0].ToAddresses).To(Equal([]string{
					"data@test.com",
				}))
				Expect(mailItems[0].Body).To(Equal("This is an email with &#39;data&#39; in the &#39;TO&#39; field"))
			})

			It("accepts and records attachments", func() {
				var err error
				var mailItems []mailslurper.MailItem
				var attachment mailslurper.Attachment

				body := "Mail with attachments"
				attachmentBody := "Header 1,Header 2\r\nValue,Value"

				to := []string{"bob@test.com"}
				msg := []byte("Content-Type: multipart/mixed; boundary=\"==b==\"\r\n" +
					"To: bob@test.com\r\n" +
					"Subject: Mail with attachments\r\n" +
					"Date: Thu, 08 Dec 2016 23:46:05 -0600 CST\r\n" +
					"\r\n" +
					"--==b==\r\n" +
					"Content-Type: text/plain\r\n" +
					"\r\n" +
					body + "\r\n" +
					"--==b==\r\n" +
					"Content-Type: text/csv\r\n" +
					"Content-Disposition: attachment; filename=\"test.csv\"\r\n" +
					"\r\n" +
					attachmentBody + "\r\n" +
					"--==b==--")

				smtp.SendMail(address, auth, from, to, msg)
				Pause()

				search := &mailslurper.MailSearch{}
				if mailItems, err = database.GetMailCollection(0, 1, search); err != nil {
					Fail(fmt.Sprintf("Error getting mail collection from database: %s", err.Error()))
				}

				Expect(len(mailItems)).To(Equal(1))
				Expect(mailItems[0].Subject).To(Equal("Mail with attachments"))
				Expect(mailItems[0].DateSent).To(Equal("2016-12-08 23:46:05"))
				Expect(mailItems[0].ContentType).To(Equal("multipart/mixed; boundary=\"==b==\""))
				Expect(mailItems[0].FromAddress).To(Equal(from))
				Expect(mailItems[0].ToAddresses).To(Equal([]string{
					"bob@test.com",
				}))
				Expect(mailItems[0].Body).To(Equal(body))
				Expect(len(mailItems[0].Attachments)).To(Equal(1))

				if attachment, err = database.GetAttachment(mailItems[0].ID, mailItems[0].Attachments[0].ID); err != nil {
					Fail(fmt.Sprintf("Unable to get the attachment: %s", err.Error()))
				}

				Expect(attachment.MailID).To(Equal(mailItems[0].ID))
				Expect(attachment.Headers.ContentType).To(Equal("text/csv"))
				Expect(attachment.Headers.FileName).To(Equal("test.csv"))
				Expect(attachment.Contents).To(Equal(attachmentBody))
			})
		})

		/*
		 * HTML Tests
		 */
		Context("that is in HTML", func() {
			/* Basic HTML */
			Context("with nothing else special about it", func() {
				It("records to the database", func() {
					var err error
					var mailItems []mailslurper.MailItem

					body := "<p>This is a basic HTML email</p>"
					to := []string{"bob@test.com"}
					msg := []byte("Content-Type: multipart/mixed; boundary=\"==b==\"\r\n" +
						"To: bob@test.com\r\n" +
						"Subject: Basic HTML Test\r\n" +
						"Date: Thu, 08 Dec 2016 23:46:05 -0600 CST\r\n" +
						"\r\n" +
						"--==b==\r\n" +
						"Content-Type: text/html; charset=\"us-ascii\"\r\n" +
						"Content-Transfer-Encoding: 7bit\r\n" +
						"\r\n" +
						body + "\r\n" +
						"--==b==--\r\n")

					smtp.SendMail(address, auth, from, to, msg)
					Pause()

					search := &mailslurper.MailSearch{}
					if mailItems, err = database.GetMailCollection(0, 1, search); err != nil {
						Fail(fmt.Sprintf("Error getting mail collection from database: %s", err.Error()))
					}

					Expect(len(mailItems)).To(Equal(1))
					Expect(mailItems[0].Subject).To(Equal("Basic HTML Test"))
					Expect(mailItems[0].DateSent).To(Equal("2016-12-08 23:46:05"))
					Expect(mailItems[0].ContentType).To(Equal("multipart/mixed; boundary=\"==b==\""))
					Expect(mailItems[0].FromAddress).To(Equal(from))
					Expect(mailItems[0].ToAddresses).To(Equal([]string{
						"bob@test.com",
					}))
					Expect(mailItems[0].Body).To(Equal(body))
				})
			})
		})
	})
})
