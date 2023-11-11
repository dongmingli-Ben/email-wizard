package main

import (
	"context"
	"database/sql"
	"email-wizard/data/logger"
	"email-wizard/data/utils"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

func serve_store_events(producer *kafka.Producer, wg *sync.WaitGroup) {
	defer wg.Done()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "store-events",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("producer", "store-events"))
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"events"}, nil)

	// // Manually assign partitions to the beginning offsets
	// partitions, _ := consumer.Assignment()
	// for _, partition := range partitions {
	// 	consumer.Seek(kafka.TopicPartition{
	// 		Topic:     partition.Topic,
	// 		Partition: partition.Partition,
	// 		Offset:    kafka.OffsetBeginning,
	// 	}, 500)
	// }

	logger.Info("data-async server started", zap.String("producer", "store-events"))

	// set up goroutine pool
	sem := semaphore.NewWeighted(MAX_EVENT_WORKER_POOL_SIZE)

	// set up channel to receive batch of messages
	ch := make(chan []kafka.Message, MAX_EVENT_BATCH_SIZE)

	go batch_message_delivery(consumer, MAX_EVENT_BATCH_SIZE, ch)

	// run consumer
	ctx := context.TODO()
	for {
		messages := <-ch
		if len(messages) == 0 {
			continue
		}
		logger.Info("received batch of messages",
			zap.Int("messages num", len(messages)),
			zap.String("producer", "store-events"))
		events := make([]map[string]interface{}, 0)
		for _, message := range messages {
			var event map[string]interface{}
			if err := json.Unmarshal(message.Value, &event); err != nil {
				logger.Error("failure in parsing received message",
					zap.String("error", err.Error()),
					zap.String("producer", "store-events"))
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key:            []byte("error_store_events"),
					Value:          message.Value,
				}, nil)
				continue
			}
			events = append(events, event)
		}
		if err := sem.Acquire(ctx, 1); err != nil {
			logger.Error("failure in acquiring semaphore",
				zap.String("error", err.Error()),
				zap.String("producer", "store-events"))
			logger.Warn("dropping batch of messages because failing to acquire semaphore",
				zap.String("producer", "store-events"))
			for _, message := range messages {
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key:            []byte("error_store_events"),
					Value:          message.Value,
				}, nil)
			}
			continue
		}
		go process_events(events, producer, sem)
	}
}

func process_events(events []map[string]interface{}, producer *kafka.Producer, sem *semaphore.Weighted) {
	defer sem.Release(1)
	defer logger.LogErrorStackTrace()
	// query DB for unparsed emails
	db, err := utils.ConnectDB()
	if err != nil {
		logger.Error("failure in connecting to DB",
			zap.String("error", err.Error()),
			zap.String("producer", "store-events"))
		for _, email := range events {
			if err := produce_email_to_topic(producer, email, "errors", "error_store_events"); err != nil {
				logger.Error("failure in producing email to errors",
					zap.String("error", err.Error()),
					zap.String("producer", "store-events"))
			}
		}
		return
	}
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		logger.Error("failure in starting transaction",
			zap.String("error", err.Error()),
			zap.String("producer", "store-events"))
		return
	}
	defer txn.Rollback()
	// insert events into DB and update event ids to emails
	err = store_events(events, txn)
	if err != nil {
		logger.Error("failure in storing events",
			zap.String("error", err.Error()),
			zap.String("producer", "store-events"))
		return
	}
	if err = txn.Commit(); err != nil {
		logger.Error("failure in committing transaction",
			zap.String("error", err.Error()),
			zap.String("producer", "store-events"))
		return
	}
	logger.Info(fmt.Sprintf("successfully stored %v events", len(events)),
		zap.String("producer", "store-events"))
}

func store_events(events []map[string]interface{}, txn *sql.Tx) error {
	query := `WITH new_events AS (
				INSERT INTO events (user_id, email_id, email_address, event_content)
				VALUES `
	for i, event := range events {
		event_bytes, err := json.Marshal(event["event"])
		if err != nil {
			return err
		}
		// escape single quotes
		event_str := strings.Replace(string(event_bytes), "'", "''", -1)
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("(%v, '%s', '%s', '%s')", event["user_id"], event["email_id"], event["address"], event_str)
	}
	query += ` RETURNING email_id, email_address, event_id)
			  UPDATE emails
			  SET event_ids = emails.event_ids || map.event_ids
			  FROM (
				  SELECT email_id, email_address, array_agg(event_id) AS event_ids
				  FROM new_events
				  GROUP BY email_id, email_address
			  ) AS map
			  WHERE emails.email_id = map.email_id AND emails.email_address = map.email_address;`
	_, err := txn.Exec(query)
	if err != nil {
		logger.Info("query fail",
			zap.String("producer", "store-events"),
			zap.String("error", err.Error()),
			zap.String("query", query))
		// fmt.Println(query)
		return err
	}
	return nil
}
