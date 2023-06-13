from msal import PublicClientApplication
import os, atexit, msal
import requests

TOKEN_CACHE_PATH = (
    "D:/OneDrive - CUHK-Shenzhen/myfolder/6/email-wizard/retrieve/msal_cache.bin"
)
APP_CLIENT_ID = "34fe7958-6ad4-438e-8218-cb028e47fe40"


def retrieve_email_outlook(user_config, n_mails: int = 50):
    cache = msal.SerializableTokenCache()
    if os.path.exists(TOKEN_CACHE_PATH):
        with open(TOKEN_CACHE_PATH, "r") as f:
            cache.deserialize(f.read())

    def update_cache():
        if cache.has_state_changed:
            with open(TOKEN_CACHE_PATH, "w") as f:
                f.write(cache.serialize())

    atexit.register(update_cache)

    app = PublicClientApplication(
        APP_CLIENT_ID,
        authority="https://login.microsoftonline.com/common",
        token_cache=cache,
    )

    result = None  # It is just an initial value. Please follow instructions below.

    # We now check the cache to see
    # whether we already have some accounts that the end user already used to sign in before.
    accounts = app.get_accounts()
    accounts = list(
        filter(lambda d: d["username"] == user_config["username"], accounts)
    )

    if accounts:
        # Now let's try to find a token in cache for this account
        result = app.acquire_token_silent(["Mail.Read"], account=accounts[0])

    if not result:
        # So no suitable token exists in cache. Let's get a new one from AAD.
        result = app.acquire_token_interactive(
            scopes=["Mail.Read"],
        )
    if "access_token" in result:
        print(f'access token acquired for {user_config["username"]}.')  # Yay!
    else:
        print(result.get("error"))
        print(result.get("error_description"))
        print(result.get("correlation_id"))  # You may need this when reporting a bug

    # Replace {access_token} with the actual access token obtained earlier
    access_token = result["access_token"]

    # Set up the API request
    url = f"https://graph.microsoft.com/v1.0/me/messages"
    headers = {
        "Authorization": "Bearer " + access_token,
        "Content-Type": "application/json",
    }
    params = {
        '$orderby': 'receivedDateTime desc',
        '$top': n_mails
    }

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

    return raw_emails
