from concurrent import futures
from confluent_kafka import Consumer, Producer, OFFSET_BEGINNING
import logger
import json
from typing import List
import asyncio
import threading
import os
from threading import Thread
import time
from gen_parse import aget_parse_results


KAFKA_CONFIG = {
    "bootstrap.servers": "kafka:29092",
    "group.id": "parse",
    "auto.offset.reset": "earliest",
}
MAX_WORKERS = 10
MAX_CONCURRENT_REQUESTS_PER_WORKER = 10


class RetryCounter:
    def __init__(
        self,
        producer: Producer,
        max_retries: int = 5,
        *args,
        **kwargs,
    ):
        self.producer = producer
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
                self.producer.produce(
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
            email = json.loads(msg.value().decode("utf-8"))
            logger.info(
                f"[{self.id}]: Successfully produced email (id = {email['email_id']}, address = {email['address']}) to topic: {msg.topic()}"
            )


async def req_retry_handler(
    email: dict,
    producer: Producer,
    max_retries=3,
    sem: asyncio.BoundedSemaphore = None,
    taskid=None,
):
    retries = max_retries
    err_str = ""
    logger.info(
        f"Parsing events for email: email_id = {email['email_id']}, address = {email['address']}"
    )
    while retries > 0:
        try:
            events = await aget_parse_results(email["item"], max_retries, semaphore=sem)
            for event in events["events"]:
                item = {
                    "user_id": email["user_id"],
                    "email_id": email["email_id"],
                    "address": email["address"],
                    "event": event,
                }
                producer.produce(
                    "events",
                    key="event",
                    value=json.dumps(item),
                    on_delivery=RetryCounter(producer, max_retries, taskid=taskid),
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

    producer.produce(
        "errors",
        key="error",
        value=json.dumps({"email": email, "error": err_str}),
    )


def task_func(
    emails: List[str],
    producer: Producer,
    max_retries=3,
):
    logger.info(
        f"start processing {len(emails)} emails with thread id {threading.get_ident()}, process id {os.getpid()} ..."
    )
    try:
        emails = [json.loads(req_msg.value()) for req_msg in emails]
        # print(json.dumps(reqs, indent=2))
        sem = asyncio.BoundedSemaphore(MAX_CONCURRENT_REQUESTS_PER_WORKER)
        tasks = [
            req_retry_handler(email, producer, max_retries, sem, i)
            for i, email in enumerate(emails)
        ]
        loop = asyncio.get_event_loop()
        loop.run_until_complete(asyncio.gather(*tasks))
        logger.info(f"Finished parsing {len(emails)} emails")
    except Exception as e:
        logger.error("Error parsing emails: {}".format(str(e)))
        logger.warning("requests: {}".format([r.value() for r in emails]))
        raise e


def serve():
    logger.info(
        f"start running async server at thread id {threading.get_ident()}, process id {os.getpid()} ..."
    )

    consumer = Consumer(KAFKA_CONFIG)
    consumer.subscribe(["new_emails"])
    producer = Producer(KAFKA_CONFIG)

    def producer_periodic_poll():
        while True:
            producer.poll(0.5)
            time.sleep(0.5)

    p_thread = Thread(target=producer_periodic_poll, daemon=True)
    p_thread.start()

    def worker_init():
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        logger.debug(f"worker thread id {threading.get_ident()} initialized ...")

    logger.info("Start consuming messages ...")

    with futures.ThreadPoolExecutor(
        max_workers=MAX_WORKERS, initializer=worker_init
    ) as executor:
        while True:
            reqs = consumer.consume(
                num_messages=MAX_CONCURRENT_REQUESTS_PER_WORKER, timeout=1.0
            )
            if len(reqs) == 0:
                continue
            logger.info("Received {} requests".format(len(reqs)))
            executor.submit(task_func, reqs, producer, 5)
            logger.info("Sent {} requests to executor".format(len(reqs)))


if __name__ == "__main__":
    logger.logger_init(
        log_dir="log", name="parse-async", level="INFO", when="D", backupCount=7
    )
    serve()
