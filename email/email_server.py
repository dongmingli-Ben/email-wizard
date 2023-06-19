from concurrent import futures
import logging
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
        print(email_id, plain_text, sep="\n", end="\n\n")
    emails = {
        'items': clean_mails
    }
    return emails


class EmailHelper(pd2_grpc.EmailHelperServicer):

    def GetEmails(self, request, context):
        logging.info(f'received request {request}')
        config = json.loads(request.config)
        n_mails = request.n_mails
        emails = get_emails(config, n_mails)
        email_json = json.dumps(emails)
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
    logging.basicConfig()
    serve()
