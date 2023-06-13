import os
import json
import re
from prompts.naive_prompt import get_prompt
from utils import get_response


def clean_to_json(response: str) -> dict:
    pattern = '\{.*"events":.*\[.*\].*\}'
    matches = re.findall(pattern, response, re.DOTALL)
    matches = sorted(matches, key=lambda s: len(s), reverse=True)
    for match in matches:
        try:
            result = json.loads(match)
            return result
        except json.JSONDecodeError:
            continue
    raise RuntimeError(f'API\'s result does not contain a valid json object, returned: {response}')


def get_parse_results(email: dict, max_patience: int = 5) -> dict:
    prompt = get_prompt(email)
    while True:
        try:
            raw_result = get_response(prompt)
            result = clean_to_json(raw_result)
            break
        except Exception as e:
            if max_patience == 0:
                print('Fail to get results after', max_patience, 'trials')
                return {"event": []}
            print('Retrying after exception:', e)
            max_patience -= 1
    return result


def gen_predictions(src_path, tgt_path):
    with open(src_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    for i, email in enumerate(data):
        print('working on email', i+1, 'total', len(data))
        results = get_parse_results(email)
        email['results'] = results
    
    with open(tgt_path, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=4, ensure_ascii=False)


if __name__ == '__main__':
    gen_predictions(r'smart_parse\data\origin-email.json', r'smart_parse\predictions\test.json')