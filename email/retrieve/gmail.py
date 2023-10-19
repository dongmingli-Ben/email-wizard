from __future__ import print_function
import base64
from datetime import datetime
import re

from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from googleapiclient.discovery import build

# If modifying these scopes, delete the file token.json.
SCOPES = ['https://www.googleapis.com/auth/gmail.readonly']

def decode_base64url(s: str):
    s = s.replace('-', '+').replace('_', '/')
    while len(s) % 4 != 0:
        s += '='
    return base64.b64decode(s).decode('utf-8')

def convert_timestamp(timestamp: str) -> str:
    # Parse the Gmail API timestamp into a datetime object
    m = re.search('(.*) \(.*\)', timestamp)
    if m is None:
        parsed_time = datetime.strptime(timestamp, "%a, %d %b %Y %H:%M:%S %z")
    else:
        parsed_time = datetime.strptime(m.group(1), "%a, %d %b %Y %H:%M:%S %z")

    # Define the desired output format
    desired_output_format = "%Y-%m-%d %H:%M:%S%z"

    # Format the time in the desired output format
    return parsed_time.strftime(desired_output_format)

# Function to decode base64 and quoted-printable encoded content
def get_raw_texts(message):
    mime: str = message['mimeType']
    contents = []
    if mime.startswith('multipart') or message['body']['size'] == 0:
        for part in message['parts']:
            children = get_raw_texts(part)
            contents += children
    elif mime.startswith('text/plain'):
        content = decode_base64url(message['body']['data'])
        contents.append(content)
    return contents


def retrieve_email_gmail(user_config: dict, n_mails: int = 50):
    creds = Credentials.from_authorized_user_info(user_config['credentials'], SCOPES)

    # List the user's Gmail inbox messages
    service = build('gmail', 'v1', credentials=creds)
    results = service.users().messages().list(userId='me', labelIds=['INBOX'], 
                                              maxResults=n_mails).execute()
    messages = results.get('messages', [])
    raw_emails = []

    for message in messages:
        msg = service.users().messages().get(userId='me', id=message['id']).execute()

        # Extract subject, sender, and content
        subject = None
        sender = None
        date = None
        recipient = None

        # Iterate through the headers to find subject and sender
        for header in msg['payload']['headers']:
            if header['name'] == 'Subject':
                subject = header['value']
            elif header['name'] == 'From':
                sender = header['value']
            elif header['name'] == 'Date':
                date = convert_timestamp(header['value'])
            elif header['name'] == 'To':
                recipient = header['value']
        contents = get_raw_texts(msg['payload'])
        
        raw_email = {
            'subject': subject,
            'sender': sender,
            'date': date,
            'recipient': [recipient],
            'content': contents
        }
        raw_emails.append((msg['id'], raw_email))
    return raw_emails

if __name__ == '__main__':
    import json
    with open('./config/gmail_user_credentials.json', 'r') as f:
        credentials = json.load(f)
    res = retrieve_email_gmail({
        'credentials': credentials,
    }, 10)
    for i, raw_email in res:
        print(i, ':',  raw_email)