import os
import json
import re
from prompts.time_and_tz import get_prompt
from utils import get_response
import asyncio


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


def get_parse_results(email: dict, max_patience: int = 5, **kwargs) -> dict:
    prompt = get_prompt(email, **kwargs)
    patience = max_patience
    # import pdb; pdb.set_trace()
    while True:
        try:
            raw_result = get_response(prompt)
            result = clean_to_json(raw_result)
            break
        except Exception as e:
            if patience == 0:
                print('Fail to get results after', max_patience, 'trials')
                return {"event": []}
            print('Retrying after exception:', e)
            patience -= 1
    return result


def gen_predictions(src_path, tgt_path, n_parallel: int = 0, **kwargs):
    with open(src_path, 'r', encoding='utf-8') as f:
        data = json.load(f)

    preds = []
    
    if n_parallel == 0:
        for i, email in enumerate(data):
            print('working on email', i+1, 'total', len(data))
            results = get_parse_results(email, **kwargs)
            email['results'] = results
            preds.append(email)
    else:
        # use async
        sem = asyncio.Semaphore(n_parallel)

        async def async_get_results(i, email):
            async with sem:
                print('working on email', i+1, 'total', len(data))
                result = await asyncio.to_thread(get_parse_results, email, **kwargs)
                print('email', i+1, 'is processed, total', len(data))
                return result
            
        async def async_predict_all(data):
            tasks = [async_get_results(i, email) for i, email in enumerate(data)]
            return await asyncio.gather(*tasks)
        
        results = asyncio.run(async_predict_all(data))
        for res, email in zip(results, data):
            email['results'] = res
            preds.append(email)

    save_dir = os.path.dirname(tgt_path)
    if not os.path.exists(save_dir):
        os.makedirs(save_dir)
    
    with open(tgt_path, 'w', encoding='utf-8') as f:
        json.dump(preds, f, indent=4, ensure_ascii=False)


if __name__ == '__main__':
    # gen_predictions(r'smart_parse\data\origin-email.json', r'smart_parse\predictions\test.json')
    gen_predictions(r'smart_parse\data\origin-email.json', r'smart_parse\predictions\time_and_tz-0614-openai-update.json', 
                    n_parallel=20, timezone='Asia/Shanghai')