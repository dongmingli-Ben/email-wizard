import imaplib
from typing import List, ByteString, Tuple


def retrieve_email_bytes(
    user_config, n_mails: int = 50
) -> List[Tuple[str, ByteString]]:
    """Retrieve the last n_mails emails"""
    imap_server = imaplib.IMAP4_SSL(user_config["imap_server"])
    status, resp = imap_server.login(user_config["username"], user_config["password"])
    if status != "OK":
        print(f'Login to {user_config["username"]} failed with {resp}')

    imap_server.select("INBOX")
    status, email_ids = imap_server.search(None, "ALL")
    if status != "OK":
        print(f'Searching {user_config["username"]} emails failed.')
        return [status]
    email_ids = email_ids[0].split()  # Convert the email IDs to a list
    raw_emails = []
    for email_id in email_ids[-n_mails:]:
        status, email_data = imap_server.fetch(email_id, "(RFC822)")
        if status != "OK":
            print(f"fail to fetch email {email_id}")
        raw_emails.append((email_id.decode("utf-8"), email_data[0][1]))

    imap_server.close()
    imap_server.logout()
    return raw_emails


if __name__ == "__main__":
    import json

    raw_emails = retrieve_email_bytes(
        user_config=json.load(open("retrieve/config.json")),
        n_mails=10,
    )
    print(raw_emails[0])
