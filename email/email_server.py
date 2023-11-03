from concurrent import futures
import logger
import json

import grpc
import service.email_pb2 as pb2
import service.email_pb2_grpc as pd2_grpc

from retrieve import retrieve_email_raw_texts
from clean.cleaner import clean_email_from_raw_texts


def get_emails(user_config, n_mails: int) -> dict:
    raw_emails = retrieve_email_raw_texts(user_config, n_mails)

    clean_mails = []
    for email_id, raw_texts in raw_emails:
        plain_text = clean_email_from_raw_texts(raw_texts)
        email = {
            'email_id': email_id,
            'item': plain_text
        }
        clean_mails.append(email)
        logger.debug(f'{email_id}:\n{plain_text}\n')
    emails = {
        'items': clean_mails
    }
    return emails


class EmailHelper(pd2_grpc.EmailHelperServicer):

    def GetEmails(self, request, context):
        logger.info(f'received request {request}')
        try:
            config = json.loads(request.config)
            n_mails = request.n_mails
            emails = get_emails(config, n_mails)
            email_json = json.dumps(emails, indent=4, ensure_ascii=False)
            logger.debug(f'return email JSON:\n {email_json}')
        except Exception as e:
            logger.error(f'Error when serving request {request}')
            raise e
        return pb2.EmailReply(message=email_json)


def serve():
    port = '50051'
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pd2_grpc.add_EmailHelperServicer_to_server(EmailHelper(), server)
    server.add_insecure_port('[::]:' + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == '__main__':
    logger.logger_init(log_dir='log', name='email',
                       level='INFO', when='D', backupCount=7)
    serve()
