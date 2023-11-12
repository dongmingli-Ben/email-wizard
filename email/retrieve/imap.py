import imaplib
import aioimaplib
from typing import List, ByteString, Tuple
import logger


def retrieve_email_bytes(
    user_config, n_mails: int = 50
) -> List[Tuple[str, ByteString]]:
    """Retrieve the last n_mails emails"""
    credentials = user_config["credentials"]
    imap_server = imaplib.IMAP4_SSL(credentials["imap_server"])
    status, resp = imap_server.login(user_config["username"], credentials["password"])
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
        status, email_data = imap_server.fetch(email_id.decode("utf-8"), "(RFC822)")
        if status != "OK":
            print(f"fail to fetch email {email_id}")
        raw_emails.append((email_id.decode("utf-8"), email_data[0][1]))

    imap_server.close()
    imap_server.logout()
    return raw_emails


async def aretrieve_email_bytes(
    user_config, n_mails: int = 50
) -> List[Tuple[str, ByteString]]:
    """Retrieve the last n_mails emails"""
    credentials = user_config["credentials"]
    imap_server = aioimaplib.IMAP4_SSL(credentials["imap_server"])
    await imap_server.wait_hello_from_server()
    status, resp = await imap_server.login(
        user_config["username"], credentials["password"]
    )
    if status != "OK":
        logger.warning(f'Login to {user_config["username"]} failed with {resp}')

    await imap_server.select("INBOX")
    status, email_ids = await imap_server.search("ALL")
    if status != "OK":
        logger.warning(f'Searching {user_config["username"]} emails failed.')
        raise RuntimeError(f'Searching {user_config["username"]} emails failed.')
    email_ids = email_ids[0].split()  # Convert the email IDs to a list
    raw_emails = []
    for email_id in email_ids[-n_mails:]:
        status, email_data = await imap_server.fetch(
            email_id.decode("utf-8"), "(RFC822)"
        )
        if status != "OK":
            logger.warning(f"fail to fetch email {email_id}")
        raw_emails.append((email_id.decode("utf-8"), bytes(email_data[1])))

    await imap_server.close()
    await imap_server.logout()
    return raw_emails


if __name__ == "__main__":
    import json
    import asyncio

    raw_emails = retrieve_email_bytes(
        user_config=json.load(open("config/imap_credentials.json")),
        n_mails=10,
    )
    print(raw_emails[0][0], raw_emails[0][1][:100])
    raw_emails = asyncio.run(
        aretrieve_email_bytes(
            user_config=json.load(open("config/imap_credentials.json")),
            n_mails=10,
        )
    )
    print(raw_emails[0][0], raw_emails[0][1][:100])
