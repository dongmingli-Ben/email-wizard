from __future__ import print_function

import logging
import json

import grpc
import service.email_pb2 as pb2
import service.email_pb2_grpc as pb2_grpc


def run(user_config, n_mails):
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    print("Will try to request the server ...")
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = pb2_grpc.EmailHelperStub(channel)
        response = stub.GetEmails(pb2.EmailRequest(config=json.dumps(user_config), n_mails=n_mails))
        print("GetEmails client received: ")
        emails = json.loads(response.message)
        print(json.dumps(emails, ensure_ascii=False, indent=4))


if __name__ == '__main__':
    logging.basicConfig()

    with open('config/outlook.json', 'r') as f:
        config = json.load(f)
    run(config, 5)
