from typing import List, Tuple
import numpy as np
import os
import datetime
from pytz import timezone
from multiprocessing.pool import ThreadPool
import asyncio
import json
from utils import get_response

class Evaluator:

    def __init__(self, log_dir: str = 'logs') -> None:
        self.noti_results = []
        self.regi_results = []
        self.acti_results = []
        if not os.path.exists(log_dir):
            os.makedirs(log_dir)
        self.log_path = os.path.join(log_dir, datetime.datetime.now().strftime('%Y-%m-%dT%H-%M-%S') + '.log')

    def __len__(self):
        return len(self.activity_results)

    def summarize_results(self, results: List[dict]) -> dict:
        p_scores, r_scores, f_scores = [], [], []
        for res in results:
            if len(res) == 0: continue
            p_scores.append(res['precision'])
            r_scores.append(res['recall'])
            f_scores.append(res['f1'])
        info = {
            'precision': np.mean(p_scores),
            'recall': np.mean(r_scores),
            'f1': np.mean(f_scores)
        }
        return info

    def summarize(self):
        noti_info = self.summarize_results(self.noti_results)
        regi_info = self.summarize_results(self.regi_results)
        acti_info = self.summarize_results(self.acti_results)
        info = {
            'notification': noti_info,
            'registration': regi_info,
            'activity': acti_info
        }
        return info

    def increment_evaluate(self, hyp: List[dict], ref: List[dict]) -> dict:
        hyp_notification = [d for d in hyp if d['event_type'] == 'notification']
        ref_notification = [d for d in ref if d['event_type'] == 'notification']
        hyp_registration = [d for d in hyp if d['event_type'] == 'registration']
        ref_registration = [d for d in ref if d['event_type'] == 'registration']
        hyp_activity = [d for d in hyp if d['event_type'] == 'activity']
        ref_activity = [d for d in ref if d['event_type'] == 'activity']

        notification_results = self.evaluate(hyp_notification, ref_notification, 'notification')
        registration_results = self.evaluate(hyp_registration, ref_registration, 'registration')
        activity_results = self.evaluate(hyp_activity, ref_activity, 'activity')

        self.noti_results.append(notification_results)
        self.regi_results.append(registration_results)
        self.acti_results.append(activity_results)

        info = {
            'notification': notification_results,
            'registration': registration_results,
            'activity': activity_results
        }

        return info
    
    def evaluate(self, hyps: List[dict], refs: List[dict], event_type: str) -> dict:
        n_hyps = len(hyps)
        n_refs = len(refs)
        if n_hyps*n_refs == 0:
            if n_refs + n_hyps == 0:
                return {}
            else:
                return {'precision': 0, 'recall': 0, 'f1': 0}
        match_arr = np.zeros((n_hyps, n_refs))
        for i in range(n_hyps):
            for j in range(n_refs):
                is_match, info = self.check_match(hyps[i], refs[j], event_type)
                if is_match:
                    match_arr[i, j] = 1
                    
                with open(self.log_path, 'a', encoding='utf-8') as f:
                    if not is_match:
                        f.write('='*30 + 'MIS-match' + '='*30 + '\n')
                    else:
                        f.write('='*30 + 'match' + '='*30 + '\n')
                    f.write(f'hyp: \n{json.dumps(hyps[i], indent=4, ensure_ascii=False)}\n')
                    f.write(f'ref: \n{json.dumps(refs[j], indent=4, ensure_ascii=False)}\n')
                    f.write(f'reason: \n{info}\n')
                    f.write('='*60 + '\n')

        p = np.sum(np.any(match_arr, axis=1)) / n_hyps
        r = np.sum(np.any(match_arr, axis=0)) / n_refs
        f = 2*p*r / (p + r) if p + r > 0 else 0
        return {'precision': p, 'recall': r, 'f1': f}
    
    @staticmethod
    def time_match(hyp: str, ref: str):
        try:
            ref_t, ref_tz = ref.split()
            hyp_t, hyp_tz = hyp.split()
            hyp_dt = datetime.datetime.strptime(hyp_t, '%Y-%m-%dT%H:%M:%S')
            dt = timezone(hyp_tz).localize(hyp_dt).astimezone(timezone(ref_tz)).strftime('%Y-%m-%dT%H:%M:%S')
            return dt == ref_t
        except ValueError as e:
            print(f'Encounter error: {e}. Hyp: {hyp}, ref: {ref}. Fall back to default evaluation.')
            return True
    
    def check_match(self, hyp: dict, ref: dict, event_type: str, repeat: int = 1) -> Tuple[bool, str]:
        prompt = self.get_match_prompt(hyp, ref, event_type)
        if 'end_time' in ref and not self.time_match(hyp['end_time'], ref['end_time']):
            return False, f'end_time {hyp["end_time"]} and {ref["end_time"]} do not match'
        if 'start_time' in ref and not self.time_match(hyp['start_time'], ref['start_time']):
            return False, f'start_time {hyp["start_time"]} and {ref["start_time"]} do not match'
        n_matches = 0
        mismatch_info = ''
        for i in range(repeat):
            response = get_response(prompt)
            is_match, info = self.parse_match_response(response)
            n_matches += is_match
            if not is_match:
                mismatch_info = info
        return n_matches/repeat > 0.5, mismatch_info
    
    def parse_match_response(self, response: str) -> Tuple[bool, str]:
        result = response.lower().startswith('yes')
        return result, response
    
    def get_match_prompt(self, hyp: dict, ref: dict, event_type: str) -> str:
        if event_type == 'notification':
            prompt = f"""Here are two notification events represented in JSON format.
            Please check the summary field to see if the two events are just one notification
            but using different summary (or talking about the same event). You do not need to be
            precise about the summary. As long as the summaries can possibly mean the same thing
            (you don't need to be picky about the details in the summaries. For example, one summary
            may contain slightly more information than the other summary),
            you can think them to be the same. For example, "2023 SME Time Capsule Graduation Photo Shoot"
            and" Reminder to the 'Time Capsule'" can be regarded as
            two summaries on the same event. Another example is "The Lazy Professional Artist" and 
            "The Lazy Professional Artist - Lecture". Although the second one is clearly about a lecture,
            you can still consider these two summaries to be the same. 
            If they are, please directly say "yes" and say no more. 
            If they are not, please first say "no" and then explain why they are not. 
            Be concise on why they are not the same.

            The two events are as below:

            {hyp}

            {ref}
            """
            return prompt
        elif event_type == 'registration':
            prompt = f"""Here are two registration events represented in JSON format.
            Please check the end_time, venue, and summary field to see if the underlying
            events of the two events are the same. You do not need to be
            precise about the summary. As long as the summaries can possibly mean the same thing
            (you don't need to be picky about the details in the summaries. For example, one summary
            may contain slightly more information than the other summary),
            you can think them to be the same. For example, "2023 SME Time Capsule Graduation Photo Shoot"
            and" Reminder to the 'Time Capsule'" can be regarded as
            two summaries on the same event. Another example is "The Lazy Professional Artist" and 
            "The Lazy Professional Artist - Lecture". Although the second one is clearly about a lecture,
            you can still consider these two summaries to be the same. 
            However, you DO need to be careful about the end_time of the events.
            If they are, please directly say "yes" and say no more. 
            If they are not, please first say "no" and then explain why they are not. 
            Be concise on why they are not the same.

            The two events are as below:

            {hyp}

            {ref}
            """
            return prompt
        elif event_type == 'activity':
            prompt = f"""Here are two activities represented in JSON format.
            Please check the start_time, end_time, venue, and summary field to see if the underlying
            events of the two activities are the same. You do not need to be
            precise about the summary. As long as the summaries can possibly mean the same thing
            (you don't need to be picky about the details in the summaries. For example, one summary
            may contain slightly more information than the other summary),
            you can think them to be the same. For example, "2023 SME Time Capsule Graduation Photo Shoot"
            and" Reminder to the 'Time Capsule'" can be regarded as
            two summaries on the same event. Another example is "The Lazy Professional Artist" and 
            "The Lazy Professional Artist - Lecture". Although the second one is clearly about a lecture,
            you can still consider these two summaries to be the same. 
            However, you DO need to be careful about the start and end time of the events.
            If they are, please directly say "yes" and say no more. 
            If they are not, please first say "no" and then explain why they are not. 
            Be concise on why they are not the same.

            The two events are as below:

            {hyp}

            {ref}
            """
            return prompt
        
    def evaluate_predictions(self, hyps: List[dict], refs: List[dict], n_parallel: int = 0):
        """
        Main endpoint for evaluation.

        Parameters
        -----------
        n_parallel: int, default: 0
            Number of threads used for evaluation. When set to 0, does not use multithreading
            for acceleration.
        """
        if n_parallel == 0:
            for i in range(len(hyps)):
                self.increment_evaluate(hyps[i]['results']['events'], refs[i]['results']['events'])
            info = self.summarize()
            return info
        
        sem = asyncio.Semaphore(n_parallel)
        
        async def async_increment_evaluate(hyp, ref):
            async with sem:
                await asyncio.to_thread(self.increment_evaluate, hyp, ref)

        async def async_evaluate(hyps, refs):
            tasks = [async_increment_evaluate(hyps[i]['results']['events'], 
                                    refs[i]['results']['events']) for i in range(len(hyps))]
            await asyncio.gather(*tasks)

        asyncio.run(async_evaluate(hyps, refs))

        info = self.summarize()
        return info


if __name__ == "__main__":
    evaluator = Evaluator()

    import json
    with open(r'smart_parse\predictions\time_and_tz-new.json', 'r', encoding='utf-8') as f:
        hyps = json.load(f)
    with open(r'smart_parse\data\raw-dataset-tz.json', 'r', encoding='utf-8') as f:
        refs = json.load(f)

    info = evaluator.evaluate_predictions(hyps, refs, n_parallel=20)

    import pandas as pd
    df = pd.DataFrame.from_dict(info)
    print(df)