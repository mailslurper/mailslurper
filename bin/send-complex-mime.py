#!/bin/env python
# -*- coding: latin-1 -*-

import smtplib
import email
import email.MIMEText
import email.MIMEImage
import email.MIMEAudio
import email.MIMEMultipart
import mimetypes
import socket
import getpass
import os
import sys
from time import gmtime, strftime

if os.name == "posix":
    import pwd


_MIME_PREAMBLE = """
This is a MIME 1.0 compliant message.

If you are able to read this, your mail program
is not MIME compliant.  You will need an alternative
mail viewer to read this message in its entirety.
"""

COMMASPACE = ", "

mailDebugPath = os.path.join(os.getcwd(), "tmp.log")



def getaddress(username):
    """Translate a username into a list of e-mail addresses.

    Note: in many cases, there will be a single address for the username
    and the list will contain one element.  In some cases the username will
    correspond to a list of e-mail addresses.  This depends on the settings
    in the mail.config module.
    """
    if username is None:
        return []

    try:
        mailnames = config.USERNAME_ADDRESSES[username]
    except KeyError:
        mailnames = username

    if isinstance(mailnames, str):
        mailnames = [mailnames]

    addresses = []
    for name in mailnames:
        if "@" in name:
            address = name
        else:
            address = name + "@" + config.MAIL_DOMAIN
        addresses.append(address)

    return addresses

def sendHTML(smtp_server, to, subject="", text = "", html = "", image_dict={}, attachments = None, from_ = None, reply_to = None, cc = None, bcc = None):

    """Send an HTML e-mail message to one or more recipients.
       Send Text email if html argument is empty.

    Parameters:
        smtp_server - smtp server url
        to - an address or list of addresses to send to.
        subject - the subject of the message.
        text - the text of the message.
        html - the html content of the message
        image_dict - dictionary of images used by the email
        attachments - a list of files to attach (pathnames).
        from_ - the From address (optional).
        reply_to - a Reply-To address (optional)
        cc - an address of list of address to CC to.

    Note: if the from address is not specified, the function attempts to
    determine it automatically from the username of the process running the
    script.

    This function uses the email module to construct an e-mail message
    (including attachments if required); and the smtplib module to send the
    message.

    Using smtplib and email directly is possible, but not terribly easy.
    This function provides a very easy way to send e-mail messages
    programmatically, which should suffice in 99% of cases.

    Exceptions:
        See the documentation for the email module (a standard Python module),
        there are a number of exceptions that can be raised if an error occurs
        during the construction of the message.

        See the documentation for smtplib.SMTP.sendmail.  That function raises
        various exceptions if an error occurs.  An exception indicates that the
        message could not be sent to to any of the addressees.

    Return value:
        As for smtplib.SMTP.sendmail.  If the function returns normally, then
        e-mail has been accepted for at least one of the addressees.  The return
        value is the same as that of smtplib.SMTP.sendmail(); i.e. a dictionary
        containing all the addresses to whom the message could not be sent.
        The dictionary entries have the form:
            address => ( error_code, error_message )
    """

    to, cc, bcc = _list(to), _list(cc), _list(bcc)

    # Remove duplicates from To:, Cc: and Bcc (so recipients don't get message twice)
    cc = [ address for address in cc if address not in to ]
    bcc = [ address for address in bcc if address not in to and address not in cc ]

    # Create the message
 #   if not attachments:
    message_root = email.MIMEMultipart.MIMEMultipart("related")
    message_alternative = email.MIMEMultipart.MIMEMultipart("alternative")
    message_alternative.attach(email.MIMEText.MIMEText(text)) # No attachments - just create a text/plain message

    #HTML
    if html != "":
        message_alternative.attach(email.MIMEText.MIMEText(html, "html"))

    message_root.attach(message_alternative)

    #LOGO, ICONS
    if image_dict:
        for key, val in image_dict.iteritems():
            if os.path.isfile(val):
                key_name=key[1:-1]
                if html and key_name in html:
                    f = open(val, "rb")
                    message_image = email.MIMEImage.MIMEImage(f.read())
                    f.close()
                    message_image.add_header("Content-Disposition", "inline")
                    message_image.add_header("Content-ID", key)
                    message_root.attach(message_image)

    #ATTACHMENTS
    if attachments != None :
        _attach(message_root, attachments)

    # Set the standard message headers
    message_root["From"] = from_
    message_root["To"] = COMMASPACE.join(to)
    message_root["Subject"] = subject
    message_root["Date"] = email.Utils.formatdate(localtime = True)

    if cc:
        message_root["CC"] = COMMASPACE.join(cc)
    if bcc:
        message_root["BCC"] = COMMASPACE.join(bcc)
    if reply_to:
        message_root["Reply-To"] = COMMASPACE.join(_list(reply_to))

    message_root.preamble = "This is a multi-part message in MIME format."

    # Send the message
    session = None
    result = None
    time_now = strftime("%a, %d %b %Y %H:%M:%S +0000", gmtime())

    old_mask = os.umask(0)
    with os.fdopen(os.open(mailDebugPath, os.O_RDWR | os.O_CREAT, 0666), 'w') as log:
        try:
            session = smtplib.SMTP(smtp_server)
        except:
            log.write( time_now + " Could not connect to mail server\n")
            log.write ("Mail server name: " + smtp_server +"\n")
            log.write ("Email addressee " + message_root["To"] +", Subj " + subject + "\n")
        if session:
            try:
                result = session.sendmail(from_, to + cc, message_root.as_string())
            except smtplib.SMTPRecipientsRefused:
                log.write(time_now + " Could not send mail to user;\n")
                log.write ("User email = " + message_root["To"])
            except smtplib.SMTPHeloError:
                log.write(time_now + " Could not connect to mail server:\n")
            except smtplib.SMTPSenderRefused:
                log.write(time_now + "Could not send mail from user:\n")
                log.write (" User email = " + from_)
            except smtplib.SMTPDataError:
                log.write(time_now + " Data error:\n")
                log.write ("Could not send the following message:\n")
                log.write (message_root.as_string())
            except:
                log.write(time_now + " Unknown mail error:\n")
            session.quit()
    log.close()
    os.umask(old_mask)

    return result


def _attach(message, attachments):
    """Attach the attachments to the message."""

    # This function is derived from source code in the Python documentation.
    # Copyright ��� 2001-2004 Python Software Foundation; All Rights Reserved

    for path in attachments:
        if not os.path.exists(path):
            raise ValueError("Does not exist: " + path)

        if not os.path.isfile(path):
            raise ValueError("Not a file: " + path)

        # Guess the content type based on the file's extension.
        contenttype, encoding = mimetypes.guess_type(path)
        if contenttype is None or encoding is not None:
            # No idea (or file is encoded) - use default type
            contenttype = "application/octet-stream"
        maintype, subtype = contenttype.split("/", 1)

        if maintype == "text":
            f = open(path)
            try:
                # Note: we should handle calculating the charset
                attachment = email.MIMEText.MIMEText(f.read(), _subtype = subtype)
            finally:
                f.close()
        else:
            f = open(path, "rb")         # open in binary mode
            try:
                # Image and audio -- use the appropriate MIME type
                if maintype == "image":
                    attachment = email.MIMEImage.MIMEImage(f.read(), _subtype = subtype)
                elif maintype == "audio":
                    attachment = email.MIMEAudio.MIMEAudio(f.read(), _subtype = subtype)
                else:
                    # Other types... encode with Base64
                    attachment = email.MIMEBase.MIMEBase(maintype, subtype)
                    attachment.set_payload(f.read())
                    email.Encoders.encode_base64(attachment)
            finally:
                f.close()

        # Set the filename parameter
        attachment.add_header("Content-Disposition", "attachment", filename = os.path.basename(path))
        message.attach(attachment)


def _get_from_address():
    """Return a From address and Reply-To address appropriate for the current user."""

    # Get the username of the current user, and their e-mail address
    username = getpass.getuser()

    addresses = getaddress(username)

    # Note: some usernames (mainly those for generic users like build and cm) do not
    # have an e-mail address.  Messages for those users goes to several different
    # addresses (i.e. get_address returns a list).

    # Single address: use it as the from address (with no Reply-To)
    if len(addresses) == 1:
        from_ = addresses[0]
        reply_to = addresses[0]

    # Multiple addresses: use username@machine as From and the
    # list of e-mail addresses as the Reply-To.
    else:
        host = socket.gethostname()
        from_ = username + "@" + config.MAIL_DOMAIN
        reply_to = addresses

    return (from_, reply_to)



def _list(v):
    """Convert a value (list, string or None) into a list."""
    if v is None:
        return []
    elif hasattr(v,'lower'):
        return [v]
    assert(isinstance(v, list))
    return v


if __name__ == "__main__":

    to = "mail@example.com"
    from_ = "mail@example.com"
    subject = "Test... "
    text = "Test message... please ignore."
    html = "<p>"+ text + "</p>"
    smtp_server = "localhost:2500"

    sendHTML(smtp_server, to, subject, text , html, {}, None, from_)


