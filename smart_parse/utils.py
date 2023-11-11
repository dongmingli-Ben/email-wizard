import logging
from openai import OpenAI, AsyncOpenAI
import os
from asyncio import BoundedSemaphore
import asyncio
import time

RETRY_WAIT_SECONDS = 5


def get_response(input: str, max_retry: int = 5) -> str:
    error = None
    client = OpenAI()
    for i in range(max_retry):
        try:
            response = client.chat.completions.create(
                model="gpt-3.5-turbo-1106",
                messages=[
                    {"role": "system", "content": "You are a helpful assistant."},
                    {"role": "user", "content": input},
                ],
                stream=True,
                response_format={"type": "json_object"},
            )
            collected_messages = []
            for chunk in response:
                chunk_message = chunk.choices[0].delta.content or ""
                collected_messages.append(chunk_message)
            message = "".join(collected_messages)
            return message
        except Exception as e:
            logging.warning(
                f"get error from openai api: {e} remaining retry times {max_retry-i-1}"
            )
            error = e
            time.sleep(RETRY_WAIT_SECONDS)
            continue
    raise error


async def aget_response(
    input: str, max_retry: int = 5, *, semaphore: BoundedSemaphore = None
):
    """
    Async version of get_response. If semaphore is provided, it will be used to control the concurrency.

    NOTE: the semaphore must be from asyncio, not from threading. Otherwise, there may be deadlocks.
    """
    error = None
    if semaphore:
        await semaphore.acquire()
    client = AsyncOpenAI()
    for i in range(max_retry):
        try:
            response = await client.chat.completions.create(
                model="gpt-3.5-turbo-1106",
                messages=[
                    {"role": "system", "content": "You are a helpful assistant."},
                    {"role": "user", "content": input},
                ],
                stream=True,
                response_format={"type": "json_object"},
            )
            collected_messages = []
            async for chunk in response:
                chunk_message = chunk.choices[0].delta.content or ""
                collected_messages.append(chunk_message)
            message = "".join(collected_messages)
            if semaphore:
                semaphore.release()
            return message
        except Exception as e:
            logging.warning(
                f"get error from openai api: {e} remaining retry times {max_retry-i-1}"
            )
            error = e
            await asyncio.sleep(RETRY_WAIT_SECONDS)
            continue
    if semaphore:
        semaphore.release()
    raise error


if __name__ == "__main__":
    sem = BoundedSemaphore(2)

    resp = get_response(
        "Please explain what is the max flow problem and how to solve it? Please wrap your response in JSON format."
    )
    print(resp)
    loop = asyncio.get_event_loop()
    resp = loop.run_until_complete(
        asyncio.gather(
            aget_response(
                "Please explain what is the max flow problem and how to solve it? Please wrap your response in JSON format.",
                semaphore=sem,
            ),
            aget_response(
                "Please explain what is the max flow problem and how to solve it? Please wrap your response in JSON format.",
                semaphore=sem,
            ),
            aget_response(
                "Please explain what is the max flow problem and how to solve it? Please wrap your response in JSON format.",
                semaphore=sem,
            ),
        )
    )
    for r in resp:
        print(r)
