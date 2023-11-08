package utils

import (
	"email-wizard/backend/logger"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func UpdateUserEventsForAccountAsync(user_id int, account map[string]interface{}) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Info(fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition))
				} else {
					logger.Info(fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition))
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "requests"
	req := map[string]interface{}{
		"config": account,
		"n_mails": N_EMAIL_RETREIVAL,
	}
	req_bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          req_bytes,
		}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	return nil
}