from concurrent import futures
import logger
import json

import grpc
import service.parse_pb2 as pb2
import service.parse_pb2_grpc as pd2_grpc

from gen_parse import get_parse_results


class Parser(pd2_grpc.ParserServicer):

    def ParseEmail(self, request, context):
        logger.info(f'received request {request}')
        try:
            email = json.loads(request.email)
            additional_info = json.loads(request.additional_info)
            result = get_parse_results(email, **additional_info)
            events_json = json.dumps(result)
            logger.debug(
                f'events: {json.dumps(result, indent=4, ensure_ascii=False)}')
        except Exception as e:
            logger.error(
                f'Uncaught error when processing for request {request}')
            raise e
        return pb2.EmailParseReply(message=events_json)


def serve():
    port = '50052'
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pd2_grpc.add_ParserServicer_to_server(Parser(), server)
    server.add_insecure_port('[::]:' + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == '__main__':
    logger.logger_init(name='parse', level='INFO')
    serve()
