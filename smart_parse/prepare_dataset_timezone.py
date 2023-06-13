import datetime
from pytz import timezone
import json
import os


def convert_time(time: str):
    format_str = '%Y-%m-%dT%H:%M:%S'
    # utc = timezone('UTC')
    # local = timezone('Asia/Shanghai')
    if time == 'unspecified':
        return time
    try:
        dt = datetime.datetime.strptime(time, format_str)
        # dt = local.localize(dt).astimezone(utc)
        return time + ' Asia/Shanghai'
    except:
        with open('tmp.txt', 'w', encoding='utf-8') as f:
            f.write(time)
        input('please open tmp.txt and correct it manually. Press enter to continue: ')
        with open('tmp.txt', 'r', encoding='utf-8') as f:
            time = f.read()
        # dt = datetime.datetime.strptime(time, format_str)
        # dt = local.localize(dt).astimezone(utc)
        return time



def standardize_dataset_time(src_path, tgt_path):
    with open(src_path, 'r', encoding='utf-8') as f:
        src = json.load(f)
    
    data = []
    for i, info in enumerate(src):
        for event in info['results']['events']:
            if event['event_type'] == 'notification':
                pass
            elif event['event_type'] == 'registration':
                event['end_time'] = convert_time(event['end_time'])
            elif event['event_type'] == 'activity':
                event['start_time'] = convert_time(event['start_time'])
                event['end_time'] = convert_time(event['end_time'])
        data.append(info)
    
    save_dir = os.path.dirname(tgt_path)
    if not os.path.exists(save_dir):
        os.makedirs(save_dir)
    with open(tgt_path, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=4, ensure_ascii=False)


if __name__ == '__main__':
    # standardize_dataset_time(
    #     r'smart_parse\data\raw-dataset.json',
    #     r'smart_parse\data\raw-dataset-tz.json'
    # )
    standardize_dataset_time(
        r'smart_parse\predictions\test.json',
        r'smart_parse\predictions\test-tz.json'
    )