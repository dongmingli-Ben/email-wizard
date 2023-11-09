import requests
import aiohttp

APP_CLIENT_ID = "34fe7958-6ad4-438e-8218-cb028e47fe40"


def retrieve_email_outlook(user_config, n_mails: int = 50):
    access_token = user_config["credentials"]["auth_token"]

    # Set up the API request
    url = f"https://graph.microsoft.com/v1.0/me/messages"
    headers = {
        "Authorization": "Bearer " + access_token,
        "Content-Type": "application/json",
    }
    params = {"$orderby": "receivedDateTime desc", "$top": n_mails}

    # Make the API request
    response = requests.get(url, params=params, headers=headers)

    # Handle the response
    raw_emails = []
    if response.status_code == 200:
        emails = response.json().get("value", [])
        for email in emails:
            # Access email properties (e.g., email['subject'], email['sender'], etc.)
            email_info = {
                "subject": email["subject"],
                "sender": email["sender"]["emailAddress"]["address"],
                "date": email["sentDateTime"],
                "recipient": [
                    e["emailAddress"]["address"] for e in email["toRecipients"]
                ],
            }
            email_info["content"] = [email["body"]["content"]]
            raw_emails.append((email["id"], email_info))
    else:
        print("Error:", response.text)
        raise RuntimeError(f"Error: fail to fetch email from graph API")

    return raw_emails


async def aretrieve_email_outlook(user_config, n_mails: int = 50):
    access_token = user_config["credentials"]["auth_token"]

    # Set up the API request
    url = f"https://graph.microsoft.com/v1.0/me/messages"
    headers = {
        "Authorization": "Bearer " + access_token,
        "Content-Type": "application/json",
    }
    params = {"$orderby": "receivedDateTime desc", "$top": n_mails}

    # Make the API request
    async with aiohttp.ClientSession() as session:
        async with session.get(url, params=params, headers=headers) as response:
            # Handle the response
            raw_emails = []
            if response.status == 200:
                emails = (await response.json()).get("value", [])
                for email in emails:
                    # Access email properties (e.g., email['subject'], email['sender'], etc.)
                    email_info = {
                        "subject": email["subject"],
                        "sender": email["sender"]["emailAddress"]["address"],
                        "date": email["sentDateTime"],
                        "recipient": [
                            e["emailAddress"]["address"] for e in email["toRecipients"]
                        ],
                    }
                    email_info["content"] = [email["body"]["content"]]
                    raw_emails.append((email["id"], email_info))
            else:
                print("Error:", response.text)

    return raw_emails


if __name__ == "__main__":
    import json
    import asyncio

    # raw_emails = retrieve_email_outlook(
    #     user_config=json.load(open("./config/outlook.json")),
    #     n_mails=10,
    # )
    # print(raw_emails[0])

    raw_emails = asyncio.run(
        aretrieve_email_outlook(
            user_config=json.load(open("./config/outlook.json")),
            n_mails=10,
        )
    )
    print(raw_emails[0])
