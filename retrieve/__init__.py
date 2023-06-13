from .utils import bytes_to_texts
from .imap import retrieve_email_bytes as retrieve_email_imap
from .pop3 import retrieve_email_bytes as retrieve_email_pop3
from .outlook import retrieve_email_outlook


def retrieve_email_bytes_list(user_config, n_mails: int = 50):
    if user_config["protocol"] == "IMAP":
        return retrieve_email_imap(user_config, n_mails)
    elif user_config["protocol"] == "POP3":
        return retrieve_email_pop3(user_config, n_mails)
    raise ValueError(f'{user_config["protocal"]} is not recognized')


def retrieve_email_raw_texts(user_config, n_mails: int = 50):
    if user_config["protocol"] == "IMAP" or user_config["protocol"] == "POP3":
        raw_emails_bytes = retrieve_email_bytes_list(user_config, n_mails)

        emails = []
        for email_id, email_bytes in raw_emails_bytes:
            email_raw_texts = bytes_to_texts(email_bytes)
            emails.append((email_id, email_raw_texts))
    elif user_config["username"].endswith("outlook.com"):
        emails = retrieve_email_outlook(user_config, n_mails)
    else:
        ValueError(f"config {user_config} is not recognized")
    return emails


if __name__ == "__main__":
    import json

    raw_emails = retrieve_email_bytes_list(
        user_config=json.load(open("config.json")),
        n_mails=10,
    )
    print(raw_emails[0])
