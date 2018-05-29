# coding=UTF-8
#
# Use this script to quickly send a bunch of mails. Useful for testing.
#
import os
import sys
import json
import time
import smtplib
import datetime
import random
from datetime import timedelta

from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email import Encoders

fromAddresses = (
    "adam@adampresley.com",
    "test@adampresley.com",
    "test@gmail.com",
    "fabio@yahoo.com",
)

toAddresses = (
    "adam@adampresley.com",
    "recipient1@gmail.com",
    "recipient2@gmail.com",
    "test@altavista.com",
    "data@test.com",
)

DATE_FORMAT_1 = "%a, %d %b %Y %H:%M:%S -0700 (UTC)"
DATE_FORMAT_2 = "%d %b %Y %H:%M:%S -0800"
DATE_FORMAT_3 = "%-d %b %Y %H:%M:%S -0800"
DATE_FORMAT_4 = "%a, %d %b %Y %H:%M:%S -0700"
DATE_FORMAT_5 = "%a, %d %b %Y %H:%M:%S -0700 UTC"
DATE_FORMAT_6 = "%a, %-d %b %Y %H:%M:%S -0700 (UTC)"
DATE_FORMAT_7 = "%a, %-d %b %Y %H:%M:%S -0700"

useSSL = False
address = "127.0.0.1"
smtpPort = 2500


def makeHTMLMessage(subject, date, dateFormat, body):
    msg = MIMEMultipart()
    html = MIMEText(body, "html")

    msg["Subject"] = subject
    msg["From"] = getRandomFrom()
    msg["To"] = getRandomTo()
    msg["Date"] = date.strftime(dateFormat)

    msg.attach(html)
    return msg


def makeTextMessage(subject, date, dateFormat, body, multipart=False):
    if multipart:
        msg = MIMEMultipart()
        msg.attach(MIMEText(body))
    else:
        msg = MIMEText(body)

    msg["Subject"] = subject
    msg["From"] = getRandomFrom()
    msg["To"] = getRandomTo()
    msg["Date"] = date.strftime(dateFormat)

    return msg


def makeMultipartMessage(subject, date, dateFormat, textBody, htmlBody):
    msg = MIMEMultipart()
    html = MIMEText(htmlBody, "html")
    text = MIMEText(textBody)

    msg["Subject"] = subject
    msg["From"] = getRandomFrom()
    msg["To"] = getRandomTo()
    msg["Date"] = date.strftime(dateFormat)

    msg.attach(text)
    msg.attach(html)
    return msg


def addAttachment(msg, filename, contentType, base64Encode=True):
    contentTypeSplit = contentType.split("/")

    part = MIMEBase(contentTypeSplit[0], contentTypeSplit[1])
    part.set_payload(open(filename, "rb").read())
    Encoders.encode_base64(part)
    #part.add_header("Content-Type", contentType)
    part.add_header("Content-Disposition",
                    "attachment; filename=\"{0}\"".format(os.path.basename(filename)))

    msg.attach(part)
    return msg


def getRandomFrom():
    return fromAddresses[random.randint(0, len(fromAddresses) - 1)]


def getRandomTo():
    return toAddresses[random.randint(0, len(toAddresses) - 1)]


def sendMail(msg):
    if not useSSL:
        server = smtplib.SMTP("{0}:{1}".format(address, smtpPort))
    else:
        server = smtplib.SMTP_SSL("{0}:{1}".format(address, smtpPort))

    fromAddress = msg["From"]
    to = [msg["To"]]

    server.sendmail(fromAddress, to, msg.as_string())
    server.quit()

#
# Seed the random generator
#
random.seed(datetime.datetime.now())

try:
    #
    # Send html with "data" in "to". This is to ensure parsing data blocks
    # do not fail.
    #
    msg = makeHTMLMessage(
        "Weird TO Address",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "<p>This is an email sent to an address with 'data' in the TO field.</p>"
    )

    sendMail(msg)

    #
    # Send plain text email
    #
    msg = makeTextMessage(
        "Plain Text Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is a plain text email.\n\nSincerely,\nAdam Presley"
    )

    sendMail(msg)

    #
    # Send plain text email with special characters in subject
    #
    msg = makeTextMessage(
        "Plain Text Email with special characters (á, é, í, ó, ú)",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is a plain text email with special characters in the subject.\n(á, é, í, ó, ú)\n\nSincerely,\nAdam Presley"
    )

    sendMail(msg)

    #
    # Send text+attachment
    #
    msg = makeTextMessage(
        "Text + Attachment Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "Plain text email with an attachment.",
        multipart=True
    )

    msg = addAttachment(msg, "../logo/logo.png", "image/png")
    sendMail(msg)

    #
    # Send html+attachment
    #
    msg = makeMultipartMessage(
        "HTML + Attachment Email activité",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with an attachment. It's got logs of >great text< & special characters.",
        "<p>This is a <strong>HTML</strong> email with an attachment. It's got lots of >great text< & special` characters.</p><p>Ceci est mon contenu accentué éàçè</p>"
    )

    msg = addAttachment(msg, "../logo/logo.png", "image/png")
    msg = addAttachment(msg, "../logo/logo.png", "image/png")

    sendMail(msg)

    #
    # Send html+CSV attachment
    #
    msg = makeMultipartMessage(
        "HTML + CSV Attachment Email Start !@#$%^&*()<>\'\"\"\' End",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with a CSV attachment.",
        "<p>This is a <strong>HTML</strong> email with a CSV attachment.</p>"
    )

    msg = addAttachment(msg, "../assets/testcsv.csv",
                        "application/octet-stream", base64Encode=False)
    sendMail(msg)

    #
    # Send html with XSS
    #
    msg = makeMultipartMessage(
        "HTML Email with XSS",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "",
        "<p>This is a <strong>HTML</strong> email with XSS stuff</p><script>alert('gotcha!');</script>"
    )

    sendMail(msg)

    #
    # Send html+JSON attachment
    #
    msg = makeMultipartMessage(
        "HTML + JSON Attachment Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with a JSON attachment.",
        "<p>This is a <strong>HTML</strong> email with a CSV attachment.</p>"
    )

    msg = addAttachment(
        msg, "../cmd/mailslurper/config.json", "application/json")
    sendMail(msg)

    #
    # Send HTML with various date formats
    #
    msg = makeMultipartMessage(
        "HTML Email with Date Format 2",
        datetime.datetime.now(),
        DATE_FORMAT_2,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_2)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 3",
        datetime.datetime.now(),
        DATE_FORMAT_3,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_3)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 4",
        datetime.datetime.now(),
        DATE_FORMAT_4,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_4)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 5",
        datetime.datetime.now(),
        DATE_FORMAT_5,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_5)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 6",
        datetime.datetime.now(),
        DATE_FORMAT_6,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_6)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 7",
        datetime.datetime.now() - timedelta(days=1),
        DATE_FORMAT_7,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_7)
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "",
        datetime.datetime.now() - timedelta(days=1),
        DATE_FORMAT_7,
        "",
        "<p>This is an email with no subject</p>"
    )

    sendMail(msg)


except Exception as e:
    print("An error occurred while trying to connect and send the email: {0}".format(
        e.message))
    print(sys.exc_info())
