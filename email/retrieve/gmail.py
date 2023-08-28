from __future__ import print_function
import base64

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
    creds = None
    token = None
    
    if user_config.get('token', {}):
        creds = Credentials.from_authorized_user_info(user_config['token'], SCOPES)
    # If there are no (valid) credentials available, let the user log in.
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            raise RuntimeError('Invalid credentials: please sign in first.')
        token = creds.to_json()

    # List the user's Gmail inbox messages
    service = build('gmail', 'v1', credentials=creds)
    results = service.users().messages().list(userId='me', labelIds=['INBOX'], 
                                              maxResults=n_mails).execute()
    messages = results.get('messages', [])
    raw_emails = []

    if not messages:
        print('No messages found.')
    else:
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
                    date = header['value']
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
    return raw_emails, token

if __name__ == '__main__':
    import json
    with open('./config/gmail_credentials.json', 'r') as f:
        credentials = json.load(f)
    with open('./config/gmail_token.json', 'r') as f:
        token = json.load(f)
    res, ret_data = retrieve_email_gmail({
        'credentials': credentials,
        'token': token
    }, 10)
    for i, raw_email in res:
        print(i, ':',  raw_email)