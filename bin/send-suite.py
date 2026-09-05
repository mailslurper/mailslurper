#!/usr/bin/env python3
# coding=UTF-8
#
# Use this script to quickly send a bunch of mails. Useful for testing.
#
import datetime
import os
import random
import smtplib
import sys
import tempfile
from datetime import timedelta
from email import encoders as Encoders
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

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
smtpPort = 1025

_HERE = os.path.dirname(os.path.abspath(__file__))

# 1x1 PNG used as a stand-in attachment so this script does not depend on
# repo artwork.
_MINI_PNG = (
    b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01"
    b"\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDATx\x9cc`\x00\x00\x00"
    b"\x02\x00\x01\xe5'\xde\xfc\x00\x00\x00\x00IEND\xaeB`\x82"
)

_SAMPLE_CSV = (
    '"id","name","title"\n'
    '1,"Adam Presley","Software Architect I"\n'
    '2,"Bob Hope","Director of Awesome"\n'
    '3,"Metallica","Speed Metal"\n'
    '4,"George Washington","President"\n'
)


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


def sampleAttachmentPath():
    fd, path = tempfile.mkstemp(suffix=".png")
    os.write(fd, _MINI_PNG)
    os.close(fd)
    return path


def sampleCSVPath():
    fd, path = tempfile.mkstemp(suffix=".csv")
    os.write(fd, _SAMPLE_CSV.encode("utf-8"))
    os.close(fd)
    return path


def addAttachment(msg, filename, contentType, base64Encode=True):
    contentTypeSplit = contentType.split("/")

    part = MIMEBase(contentTypeSplit[0], contentTypeSplit[1])
    part.set_payload(open(filename, "rb").read())
    Encoders.encode_base64(part)
    part.add_header(
        "Content-Disposition",
        'attachment; filename="{0}"'.format(os.path.basename(filename)),
    )

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


random.seed(datetime.datetime.now().timestamp())

try:
    msg = makeHTMLMessage(
        "Weird TO Address",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "<p>This is an email sent to an address with 'data' in the TO field.</p>",
    )

    sendMail(msg)

    msg = makeTextMessage(
        "Plain Text Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is a plain text email.\n\nSincerely,\nAdam Presley",
    )

    sendMail(msg)

    msg = makeTextMessage(
        "Plain Text Email with special characters (á, é, í, ó, ú)",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is a plain text email with special characters in the subject.\n(á, é, í, ó, ú)\n\nSincerely,\nAdam Presley",
    )

    sendMail(msg)

    msg = makeTextMessage(
        "Text + Attachment Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "Plain text email with an attachment.",
        multipart=True,
    )

    sample = sampleAttachmentPath()
    msg = addAttachment(msg, sample, "image/png")
    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML + Attachment Email activité",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with an attachment. It's got logs of >great text< & special characters.",
        "<p>This is a <strong>HTML</strong> email with an attachment. It's got lots of >great text< & special` characters.</p><p>Ceci est mon contenu accentué éàçè</p>",
    )

    msg = addAttachment(msg, sample, "image/png")
    msg = addAttachment(msg, sample, "image/png")

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML + CSV Attachment Email Start !@#$%^&*()<>'\"\"' End",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with a CSV attachment.",
        "<p>This is a <strong>HTML</strong> email with a CSV attachment.</p>",
    )

    csv_path = sampleCSVPath()
    msg = addAttachment(msg, csv_path, "application/octet-stream", base64Encode=False)
    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with XSS",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "",
        "<p>This is a <strong>HTML</strong> email with XSS stuff</p><script>alert('gotcha!');</script>",
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML + JSON Attachment Email",
        datetime.datetime.now(),
        DATE_FORMAT_1,
        "This is an HTML email with a JSON attachment.",
        "<p>This is a <strong>HTML</strong> email with a CSV attachment.</p>",
    )

    msg = addAttachment(
        msg,
        os.path.join(_HERE, "..", "cmd", "mylslurper", "config.json"),
        "application/json",
    )
    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 2",
        datetime.datetime.now(),
        DATE_FORMAT_2,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_2
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 3",
        datetime.datetime.now(),
        DATE_FORMAT_3,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_3
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 4",
        datetime.datetime.now(),
        DATE_FORMAT_4,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_4
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 5",
        datetime.datetime.now(),
        DATE_FORMAT_5,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_5
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 6",
        datetime.datetime.now(),
        DATE_FORMAT_6,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_6
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "HTML Email with Date Format 7",
        datetime.datetime.now() - timedelta(days=1),
        DATE_FORMAT_7,
        "",
        "<p>This is an email where the date in the header is formatted with {0}</p>".format(
            DATE_FORMAT_7
        ),
    )

    sendMail(msg)

    msg = makeMultipartMessage(
        "",
        datetime.datetime.now() - timedelta(days=1),
        DATE_FORMAT_7,
        "",
        "<p>This is an email with no subject</p>",
    )

    sendMail(msg)

except Exception as e:
    print("An error occurred while trying to connect and send the email: {0}".format(e))
    print(sys.exc_info())
