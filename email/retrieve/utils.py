from typing import ByteString, Dict, List
from email import policy
from email.parser import BytesParser
from email.message import EmailMessage


def bytes_to_texts(email_bytes: ByteString) -> Dict[str, str]:
    # Create a BytesParser object and parse the raw email
    parser = BytesParser(policy=policy.default)
    message = parser.parsebytes(email_bytes)

    # Extract relevant information from the parsed email
    subject = message["Subject"]
    sender = message["From"]
    recipient = message["To"]
    date = message["Date"]

    email_info = {
        "subject": subject,
        "sender": sender,
        "date": date,
        "recipient": recipient if isinstance(recipient, list) else [recipient],
    }

    content = extract_raw_texts(message)
    email_info["content"] = content

    return email_info


def extract_raw_texts(message: EmailMessage) -> List[str]:
    texts = []
    if message.is_multipart():
        for part in message.iter_parts():
            text_list = extract_raw_texts(part)
            texts += text_list
    else:
        content_type = message.get_content_type()
        if content_type in ["text/plain", "text/html"]:
            texts.append(message.get_content())
    return texts
