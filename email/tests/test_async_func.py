import json
import asyncio
import sys
import os
import time

sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))

import logger
from retrieve import retrieve_email_raw_texts, aretrieve_email_raw_texts

logger.logger_init(name="test")
loop = asyncio.get_event_loop()

t0 = time.time()
raw_emails = retrieve_email_raw_texts(
    user_config=json.load(open("./config/outlook.json")),
    n_mails=10,
)
print(f"Sync: fetched 10 mails in {time.time() - t0:.2f}s")

t0 = time.time()
araw_emails = loop.run_until_complete(
    aretrieve_email_raw_texts(
        user_config=json.load(open("./config/outlook.json")),
        n_mails=10,
    )
)
print(f"Async: fetched 10 mails in {time.time() - t0:.2f}s")


print("Async == Sync: ", str(raw_emails) == str(araw_emails))


async def fetch_with_retry(func, max_retry=5, *args, **kwargs):
    retry = max_retry
    while retry > 0:
        try:
            return await func(*args, **kwargs)
        except Exception as e:
            print(f"Error: {e}, retrying...")
            retry -= 1
    raise RuntimeError(f"Error: fail to fetch email after {max_retry} retries")


t0 = time.time()
_ = loop.run_until_complete(
    asyncio.gather(
        *[
            fetch_with_retry(
                aretrieve_email_raw_texts,
                5,
                json.load(open("./config/outlook.json")),
                n_mails=10,
            )
            for _ in range(10)
        ]
    )
)
print(f"Async: fetched 100 mails in 10 calls in {time.time() - t0:.2f}s")
