import logging
import openai
import os

# Load your API key from an environment variable or secret management service
openai.api_key = os.getenv('OPENAI_API_KEY')


def get_response(input: str, max_retry: int = 5) -> str:
    error = None
    for i in range(max_retry):
        try:
            response = openai.chat.completions.create(
                model="gpt-3.5-turbo-1106",
                messages=[
                    {"role": "system", "content": "You are a helpful assistant."},
                    {"role": "user", "content": input},
                ],
                stream=True,
                response_format={"type": "json_object"}
            )
            collected_messages = []
            for chunk in response:
                chunk_message = chunk.choices[0].delta.content or ""
                collected_messages.append(chunk_message)
            message = ''.join(collected_messages)
            return message
        except Exception as e:
            logging.warning(
                f'get error from openai api: {e} remaining retry times {max_retry-i-1}')
            error = e
            continue
    raise error


if __name__ == '__main__':
    resp = get_response(
        "Please explain what is the max flow problem and how to solve it? Please wrap your response in JSON format.")
    print(resp)
