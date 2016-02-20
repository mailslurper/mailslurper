#
# Use this script to quickly send a bunch of mails. Useful for testing.
#
import sys
import json
import time
import smtplib
import urllib2
import datetime

from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email import Encoders

def getQuote():
	response = urllib2.urlopen("http://www.iheartquotes.com/api/v1/random?format=json")
	obj = json.loads(response.read())

	quoteLines = obj["quote"].split("--")

	if len(quoteLines) > 0:
		return {
			"quote": quoteLines[0].strip(),
			"source": "Unknown" if len(quoteLines) <= 1 else quoteLines[1].strip(),
		}
	else:
		return {
			"quote": "No quote",
			"source": "Adam Presley"
		}

if __name__ == "__main__":
	sendMultipartMails = True
	sendTextOnlyMails = True
	sendAttachmentMails = True

	numMails = 5
	address = "127.0.0.1"
	smtpPort = 2500

	me = "someone@another.com"
	me2 = "different@another.com"
	to = "bob@bobtestingmailslurper.com"
	to2 = "jim@bobtestingmailslurper.com"

	try:
		#
		# Send text+html emails
		#
		#if sendMultipartMails:
			# for index in range(numMails):
			# 	quote = getQuote()
			#
			# 	textBody = "Hello,\nHere is today's quote.\n\n{0}\n  -- {1}\n\nSincerely,\nAdam Presley".format(quote["quote"], quote["source"])
			# 	htmlBody = "<p>Hello,</p><p>Here is today's quote.</p><p><em>{0}</em><br />&nbsp;&nbsp;-- {1}</p><p>Sincerely,<br />Adam Presley</p>".format(quote["quote"], quote["source"],)
			#
			# 	text = MIMEText(textBody, "plain")
			# 	html = MIMEText(htmlBody, "html")
			#
			# 	msg = MIMEMultipart("alternative")
			#
			# 	msg["Subject"] = "Quote From {0}".format(quote["source"])
			# 	msg["From"] = me
			# 	msg["To"] = to
			# 	msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S +0000 UTC")
			#
			# 	msg.attach(text)
			# 	msg.attach(html)
			#
			# 	server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			# 	server.sendmail(me, [to], msg.as_string())
			# 	server.quit()

				#time.sleep(2)

		#
		# Send plain text emails
		#
		if sendTextOnlyMails:
			for index in range(numMails):
				textBody = "Hello,\nI am plain text mail #{0}.\n\nSincerely,\nAdam Presley".format(index,)

				msg = MIMEText(textBody)

				msg["Subject"] = "Text Mail #{0}".format(index,)
				msg["From"] = me
				msg["To"] = to2
				msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

				server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
				server.sendmail(me, [to2], msg.as_string())
				server.quit()

				#time.sleep(1)

		#
		# Send plain text email with no subject
		#
		if sendTextOnlyMails:
			textBody = "Hello,\nI am plain text mail with no subject.\n\nSincerely,\nAdam Presley"

			msg = MIMEText(textBody)

			msg["Subject"] = ""
			msg["From"] = me2
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me2, [to], msg.as_string())
			server.quit()

		#
		# Send text+attachment
		#
		if sendAttachmentMails:
			textBody = "Hello,\nI am plain text mail with an attachment.\n\nSincerely,\nAdam Presley"

			msg = MIMEMultipart()

			msg["Subject"] = "Text+Attachment Mail"
			msg["From"] = me
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(MIMEText(textBody))

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./MailSlurperLogo.png", "rb").read())
			Encoders.encode_base64(part)
			part.add_header("Content-Type", "image/png")
			part.add_header("Content-Disposition", "attachment; filename=\"MailSlurperLogo.png\"")
			msg.attach(part)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to], msg.as_string())
			server.quit()

			#time.sleep(1)

			#
			# Send html+attachment
			#
			htmlBody = "<p>This is a <strong>HTML</strong> email with an attachment.It's got lots of >great text< & special` characters.</p>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "Adam's HTML+Attachment Mail"
			msg["From"] = me
			msg["To"] = to2
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(html)

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./MailSlurperLogo.png", "rb").read())
			Encoders.encode_base64(part)
			part.add_header("Content-Type", "image/png")
			part.add_header("Content-Disposition", "attachment; filename=\"MailSlurperLogo1.png\"")
			msg.attach(part)

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./MailSlurperLogo.png", "rb").read())
			Encoders.encode_base64(part)
			part.add_header("Content-Type", "image/png")
			part.add_header("Content-Disposition", "attachment; filename=\"MailSlurperLogo2.png\"")
			msg.attach(part)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to2], msg.as_string())
			server.quit()

			#
			# Send html+CSV attachment
			#
			htmlBody = "<p>This is a <strong>HTML</strong> email with a CSV attachment.</p>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "HTML+CSV Attachment Mail"
			msg["From"] = me2
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(html)

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./test-files/testcsv.csv", "rb").read())
			#Encoders.encode_base64(part)
			part.add_header("Content-Type", "application/octet-stream")
			part.add_header("Content-Disposition", "attachment; filename=\"testcsv.csv\"")
			msg.attach(part)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me2, [to], msg.as_string())
			server.quit()

			#
			# Send html with XSS
			#
			htmlBody = "<p>This is a <strong>HTML</strong> email with XSS stuff</p><script>alert('gotcha!');</script>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "HTML Mail with XSS"
			msg["From"] = me
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(html)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to], msg.as_string())
			server.quit()

			#
			# Send html+attachment (JSON)
			#
			htmlBody = "<p>This is a <strong>HTML</strong> email with a JSON attachment.</p>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "HTML+JSON Attachment Mail"
			msg["From"] = me
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(html)

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./config.json", "rb").read())
			Encoders.encode_base64(part)
			part.add_header("Content-Type", "application/json")
			part.add_header("Content-Disposition", "attachment; filename=\"config.json\"")
			msg.attach(part)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to], msg.as_string())
			server.quit()

			#time.sleep(1)

			#
			# Send html+attachment with filename in content-type as "name"
			#
			htmlBody = "<p>This is a <strong>HTML</strong> email with an attachment done differently.</p>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "HTML+Attachment Mail 2"
			msg["From"] = me
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%a, %d %b %Y %H:%M:%S -0700 (UTC)")

			msg.attach(html)

			part = MIMEBase("multipart", "mixed")
			part.set_payload(open("./MailSlurperLogo.png", "rb").read())
			Encoders.encode_base64(part)
			part.add_header("Content-Type", "image/png; name=MailSlurperLogo.png")
			part.add_header("Content-Disposition", "attachment;")
			msg.attach(part)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to], msg.as_string())
			server.quit()

			#
			# Send html with 4th form date format
			#
			htmlBody = "<p>This is a <strong>HTML</strong>. This has a weird date in the header.</p>"

			msg = MIMEMultipart()
			html = MIMEText(htmlBody, "html")

			msg["Subject"] = "Adam's HTML+4th Format Date in Header"
			msg["From"] = me
			msg["To"] = to
			msg["Date"] = datetime.datetime.now().strftime("%d %b %Y %H:%M:%S -0800")

			msg.attach(html)

			server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
			server.sendmail(me, [to], msg.as_string())
			server.quit()



	except Exception as e:
		print("An error occurred while trying to connect and send the email: {0}".format(e.message))
		print(sys.exc_info())
