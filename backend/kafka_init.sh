#!/bin/bash

# List of topic names
topics=("requests" "emails" "new_emails" "events" "errors")

# Kafka broker connection information
bootstrap_server="kafka:29092"

for topic in "${topics[@]}"; do
  # Create each topic using kafka-topics.sh
  docker exec -d kafka kafka-topics.sh --create --topic "$topic" --partitions 1 --replication-factor 1 --bootstrap-server "$bootstrap_server"
done
