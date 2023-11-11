from confluent_kafka import Consumer, Producer, OFFSET_BEGINNING
import json
import sys
import os

sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))

import logger


def produce_emails(email, n):
    producer = Producer({"bootstrap.servers": "kafka:29092"})
    for i in range(n):
        producer.produce("new_emails", key=str(i), value=json.dumps(email))
    producer.flush()
    logger.info(f"Produced {n} emails")


def main():
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    logger.logger_init(log_dir="../log", name="benchmark_async", level="INFO")
    with open("example_input.json") as f:
        email = json.load(f)
        email = json.loads(email["email"])
    item = {
        "user_id": 1,
        "email_id": "id",
        "address": "address@outlook.com",
        "protocol": "outlook",
        "item": email,
    }
    produce_emails(item, 200)


if __name__ == "__main__":
    main()
