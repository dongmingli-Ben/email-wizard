import poplib
from typing import List, ByteString, Tuple
import asyncio


def retrieve_email_bytes(
    user_config, n_mails: int = 50
) -> List[Tuple[str, ByteString]]:
    """Retrieve the last n_mails emails"""
    credentials = user_config["credentials"]
    pop_server = poplib.POP3_SSL(credentials["pop3_server"])
    pop_server.user(user_config["username"])
    pop_server.pass_(credentials["password"])
    num_emails = len(pop_server.list()[1])

    raw_emails = []
    for i in range(max(1, num_emails - n_mails + 1), num_emails + 1):
        _, email_data, _ = pop_server.retr(i)
        raw_emails.append((str(i), email_data))

    return raw_emails


async def aretrieve_email_bytes(
    user_config, n_mails: int = 50
) -> List[Tuple[str, ByteString]]:
    """Retrieve the last n_mails emails"""
    raw_emails = await asyncio.to_thread(retrieve_email_bytes, user_config, n_mails)
    return raw_emails
