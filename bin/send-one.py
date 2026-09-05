#!/usr/bin/env python3
"""Send a single test email to a local MylSlurper SMTP listener."""

import smtplib
import sys
from datetime import datetime
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText


def send_test_email(
    smtp_host="localhost",
    smtp_port=1025,
    from_addr="test@example.com",
    to_addr="recipient@example.com",
    subject="Test Email from Python Script",
):
    body = f"""Hello from MylSlurper Test!

This is a test email sent from a Python script at {datetime.now()}.

MylSlurper should capture this email and display it in the web interface.

Features tested:
- SMTP connection to MylSlurper
- Email capture and storage
- Web interface display

Best regards,
Python Test Script"""

    msg = MIMEMultipart()
    msg["From"] = from_addr
    msg["To"] = to_addr
    msg["Subject"] = subject
    msg.attach(MIMEText(body, "plain"))

    print(f"Sending email to MylSlurper at {smtp_host}:{smtp_port}...")
    print(f"From: {from_addr}")
    print(f"To: {to_addr}")
    print(f"Subject: {subject}")

    try:
        server = smtplib.SMTP_SSL(smtp_host, smtp_port, local_hostname="localhost")
        server.send_message(msg)
        server.quit()
        print("\nEmail sent successfully!")
        print("Check the web interface at: http://localhost:4436")
        return True
    except Exception as e:
        print(f"\nError sending email: {e}")
        return False


if __name__ == "__main__":
    to_addr = sys.argv[1] if len(sys.argv) > 1 else "recipient@example.com"
    send_test_email(to_addr=to_addr)
