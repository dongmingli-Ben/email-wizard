from concurrent import futures
from confluent_kafka import Consumer, Producer, OFFSET_BEGINNING
from endpoint import aget_emails
import logger
import json
from typing import List
import asyncio
import threading
import os
from copy import deepcopy
from threading import Thread

KAFKA_CONFIG = {
    "bootstrap.servers": "kafka:29092",
    "group.id": "email",
    "auto.offset.reset": "earliest",
}


def req_to_str_log(req: dict):
    log_req = deepcopy(req)
    log_req["config"]["credentials"] = "********"
    return json.dumps(log_req, indent=2)


class RetryCounter:
    def __init__(
        self,
        producer: Producer,
        err_producer: Producer,
        max_retries: int = 5,
        *args,
        **kwargs,
    ):
        self.producer = producer
        self.err_producer = err_producer
        self.retries = max_retries
        self.id = kwargs.get("taskid", -1)

    def __call__(self, err, msg):
        if err is not None:
            logger.warning(f"Error when producing message to topic: {err}")
            self.retries -= 1
            if self.retries > 0:
                self.producer.produce(
                    msg.topic(),
                    key=msg.key().decode("utf-8"),
                    value=msg.value().decode("utf-8"),
                    on_delivery=self,
                )
                logger.info(f"Retry producing message to topic: {msg.topic()}")
            else:
                logger.error(f"Failed to produce message to topic: {msg.topic()}")
                self.err_producer.produce(
                    "errors",
                    key="error",
                    value=json.dumps(
                        {
                            "key": msg.key().decode("utf-8"),
                            "msg": msg.value().decode("utf-8"),
                            "error": f"fail to produce message to topic: {msg.topic()}",
                            "description": str(err),
                        }
                    ),
                )
        else:
            logger.debug(
                f"[{self.id}]: Successfully produced message to topic: {msg.topic()}"
            )


async def req_retry_handler(
    req: dict,
    email_producer: Producer,
    error_producer: Producer,
    max_retries=3,
    taskid=None,
):
    retries = max_retries
    err_str = ""
    logger.info(f"Retrieving emails for request: {req_to_str_log(req)}")
    while retries > 0:
        try:
            emails = await aget_emails(req["config"], req["n_mails"])
            for email in emails["items"]:
                email_producer.produce(
                    "emails",
                    key="email",
                    value=json.dumps(email),
                    on_delivery=RetryCounter(
                        email_producer, error_producer, max_retries, taskid=taskid
                    ),
                )
            return
        except Exception as e:
            retries -= 1
            err_str = str(e)
            logger.error(
                "Error retriving emails for request: {}, retrying in 1s ...".format(
                    err_str
                )
            )
            await asyncio.sleep(1)

    error_producer.produce(
        "errors",
        key="error",
        value=json.dumps({"req": req, "error": err_str}),
    )


def task_func(
    reqs: List[str], email_producer: Producer, error_producer: Producer, max_retries=3
):
    logger.info(
        f"start processing {len(reqs)} requests with thread id {threading.get_ident()}, process id {os.getpid()} ..."
    )
    try:
        reqs = [json.loads(req_msg.value()) for req_msg in reqs]
        # print(json.dumps(reqs, indent=2))
        tasks = [
            req_retry_handler(req, email_producer, error_producer, max_retries, i)
            for i, req in enumerate(reqs)
        ]
        loop = asyncio.get_event_loop()
        loop.run_until_complete(asyncio.gather(*tasks))
        logger.info(f"Finished processing {len(reqs)} requests")
    except Exception as e:
        logger.error("Error parsing requests: {}".format(str(e)))
        logger.warning("requests: {}".format([r.value() for r in reqs]))
        raise e


def serve():
    logger.info(
        f"start running async server at thread id {threading.get_ident()}, process id {os.getpid()} ..."
    )

    consumer = Consumer(KAFKA_CONFIG)
    consumer.subscribe(["requests"])
    email_producer = Producer(KAFKA_CONFIG)
    error_producer = Producer(KAFKA_CONFIG)

    def producer_periodic_poll():
        while True:
            email_producer.poll(0.5)
            error_producer.poll(0.5)

    p_thread = Thread(target=producer_periodic_poll, daemon=True)
    p_thread.start()

    def worker_init():
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        logger.debug(f"worker thread id {threading.get_ident()} initialized ...")

    with futures.ThreadPoolExecutor(
        max_workers=10, initializer=worker_init
    ) as executor:
        while True:
            reqs = consumer.consume(num_messages=10, timeout=1.0)
            if len(reqs) == 0:
                continue
            logger.info("Received {} requests".format(len(reqs)))
            executor.submit(task_func, reqs, email_producer, error_producer, 5)
            logger.info("Sent {} requests to executor".format(len(reqs)))


if __name__ == "__main__":
    logger.logger_init(
        log_dir="log", name="email-async", level="INFO", when="D", backupCount=7
    )
    serve()
