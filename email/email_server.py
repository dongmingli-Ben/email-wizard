from concurrent import futures
from endpoint import get_emails
import logger
import json

import grpc
import service.email_pb2 as pb2
import service.email_pb2_grpc as pd2_grpc


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
