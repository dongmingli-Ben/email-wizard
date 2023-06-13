import argparse
import json
import os
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
    parser.add_argument(
        "--output",
        type=str,
        default="data/origin-email.json",
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
        plain_text_info = clean_email_from_raw_texts(raw_texts)
        clean_mails.append(plain_text_info)
        
    save_dir = os.path.dirname(args.output)
    if not os.path.exists(save_dir):
        os.makedirs(save_dir)
    with open(args.output, "w", encoding='utf-8') as f:
        json.dump(clean_mails, f, indent=4, ensure_ascii=False)


if __name__ == "__main__":
    main()
