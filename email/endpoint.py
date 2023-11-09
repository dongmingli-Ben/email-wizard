import logger

from retrieve import aretrieve_email_raw_texts, retrieve_email_raw_texts
from clean.cleaner import clean_email_from_raw_texts


def get_emails(user_config, n_mails: int) -> dict:
    raw_emails = retrieve_email_raw_texts(user_config, n_mails)

    clean_mails = []
    for email_id, raw_texts in raw_emails:
        plain_text = clean_email_from_raw_texts(raw_texts)
        email = {"email_id": email_id, "address": user_config['username'], "item": plain_text}
        clean_mails.append(email)
        logger.debug(f"{email_id}:\n{plain_text}\n")
    emails = {"items": clean_mails}
    return emails


async def aget_emails(user_config, n_mails: int) -> dict:
    """
    Retrieve emails from user's email account and clean them.
    
    Return:
        * emails: `dict`

        >>> {"items": 
        >>>     [
        >>>         {
        >>>             "email_id": str, 
        >>>             "item": 
        >>>                 {
        >>>                     "subject": str, ..., 
        >>>                     "content": str
        >>>                 }
        >>>         }
        >>>     ]
        >>> }
    """
    raw_emails = await aretrieve_email_raw_texts(user_config, n_mails)

    clean_mails = []
    for email_id, raw_texts in raw_emails:
        plain_text = clean_email_from_raw_texts(raw_texts)
        email = {"email_id": email_id, "address": user_config['username'], "item": plain_text}
        clean_mails.append(email)
        logger.debug(f"{email_id}:\n{plain_text}\n")
    emails = {"items": clean_mails}
    return emails
