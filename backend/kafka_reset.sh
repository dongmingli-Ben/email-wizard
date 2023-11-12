#!/bin/bash

# Usage:
# bash kafka_reset.sh <topic_name>
# 
# NOTE: DO NOT RESET __consumer_offsets
# Otherise, consumers cannot receive messages

set -e

if [ $1 == "__consumer_offsets" ]; then
    echo "DO NOT RESET __consumer_offsets"
    exit 1
fi

docker exec -it \
    kafka bash -c "
        # set -e
        kafka-topics.sh \
            --delete \
            --bootstrap-server kafka:29092 \
            --topic $1
        # kafka-topics.sh \
        #     --list \
        #     --bootstrap-server kafka:29092
        kafka-topics.sh \
            --create --topic $1 \
            --partitions 1 --replication-factor 1 \
            --bootstrap-server kafka:29092
        "