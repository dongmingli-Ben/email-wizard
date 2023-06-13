import argparse
import json
from retrieve import retrieve_email_raw_texts
from clean.cleaner import clean_email_from_raw_texts


def parse_args():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "--n-mails", type=int, default=5, help="Number of emails to analyze"
    )
    parser.add_argument(
        "--config",
        type=str,
        default="retrieve/config.json",
        help="user configuration file",
    )

    args = parser.parse_args()

    return args


def main():
    args = parse_args()

    with open(args.config, "r") as f:
        user_config = json.load(f)
    raw_emails = retrieve_email_raw_texts(user_config, args.n_mails)

    clean_mails = []
    for email_id, raw_texts in raw_emails:
        plain_text = clean_email_from_raw_texts(raw_texts)
        clean_mails.append((email_id, plain_text))
        print(email_id, plain_text, sep="\n", end="\n\n")


if __name__ == "__main__":
    main()
