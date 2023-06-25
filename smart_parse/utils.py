import logging
import openai
import os

# Load your API key from an environment variable or secret management service
openai.api_key = os.getenv('OPENAI_API_KEY')


def get_response(input: str, max_retry: int = 5) -> str:
    error = None
    for i in range(max_retry):
        try:
            response = openai.ChatCompletion.create(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": "You are a helpful assistant."},
                    {"role": "user", "content": input},
                ],
                stream=True
            )
            collected_chunks = []
            collected_messages = []
            for chunk in response:
                collected_chunks.append(chunk)
                chunk_message = chunk['choices'][0]['delta']
                collected_messages.append(chunk_message)
            message = ''.join([m.get('content', '') for m in collected_messages])
            return message
        except Exception as e:
            logging.warning(f'get error from openai api: {e} remaining retry times {max_retry-i-1}')
            error = e
            continue
    raise error