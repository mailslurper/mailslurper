#!/usr/bin/env python3
"""
Simple script to send test emails to MailSlurper
"""

import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from datetime import datetime
import sys

def send_test_email(smtp_host='localhost', smtp_port=2500, 
                   from_addr='test@example.com', to_addr='recipient@example.com',
                   subject='Test Email from Python Script'):
    """Send a test email to MailSlurper"""
    
    # Email body
    body = f"""Hello from MailSlurper Test!

This is a test email sent from a Python script at {datetime.now()}.

MailSlurper should capture this email and display it in the web interface.

Features tested:
- SMTP connection to MailSlurper
- Email capture and storage
- Web interface display

Best regards,
Python Test Script"""

    # Create the email message
    msg = MIMEMultipart()
    msg['From'] = from_addr
    msg['To'] = to_addr
    msg['Subject'] = subject
    msg.attach(MIMEText(body, 'plain'))

    # Send the email
    print(f'Sending email to MailSlurper at {smtp_host}:{smtp_port}...')
    print(f'From: {from_addr}')
    print(f'To: {to_addr}')
    print(f'Subject: {subject}')

    try:
        server = smtplib.SMTP(smtp_host, smtp_port)
        server.send_message(msg)
        server.quit()
        print('\n✅ Email sent successfully!')
        print('📧 Check MailSlurper web interface at: http://localhost:8080')
        print('🔧 Service API available at: http://localhost:8085')
        return True
    except Exception as e:
        print(f'\n❌ Error sending email: {e}')
        return False

if __name__ == '__main__':
    # Allow custom parameters via command line
    if len(sys.argv) > 1:
        to_addr = sys.argv[1]
    else:
        to_addr = 'recipient@example.com'
    
    send_test_email(to_addr=to_addr)
