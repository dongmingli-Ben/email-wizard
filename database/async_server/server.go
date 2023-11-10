package main

import (
	"context"
	"email-wizard/data/logger"
	"email-wizard/data/utils"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var MAX_BATCH_SIZE int = 20
var MAX_BATCH_WAIT_TIME = time.Second
var MAX_WORKER_POOL_SIZE int64 = 2
var PRODUCER_FLUSH_INTERVAL = time.Millisecond * 500
var PRODUCER_FLUSH_WAIT_TIME int = 500
var CONSUMER_WAIT_TIME_PER_MESSAGE = time.Millisecond * 50

func produce_email_to_topic(producer *kafka.Producer, email interface{}, topic string, key string) error {
	email_json, err := json.Marshal(email)
	if err != nil {
		return err
	}
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key: []byte(key),
		Value: email_json,
	}, nil)
	return nil
}

func batch_message_delivery(consumer *kafka.Consumer, ch chan []kafka.Message) {
	t0 := time.Now()
	buffer := make([]kafka.Message, 0)
	for {
		if time.Since(t0) > MAX_BATCH_WAIT_TIME || len(buffer) >= MAX_BATCH_SIZE {
			ch <- buffer
			consumer.Commit()
			buffer = make([]kafka.Message, 0)
			t0 = time.Now()
		}
		msg, err := consumer.ReadMessage(CONSUMER_WAIT_TIME_PER_MESSAGE)
		if err != nil && err.(kafka.Error).Code() != kafka.ErrTimedOut {
			logger.Error("failure in reading message", zap.String("error", err.Error()))
			continue
		}
		if err == nil {
			buffer = append(buffer, *msg)
		}
	}
}

func process_emails(emails []map[string]interface{}, producer *kafka.Producer, sem *semaphore.Weighted) {
	defer sem.Release(1)
	defer logger.LogErrorStackTrace()
	// query DB for unparsed emails
	db, err := utils.ConnectDB()
	if err != nil {
		logger.Error("failure in connecting to DB", zap.String("error", err.Error()))
		for _, email := range emails {
			if err := produce_email_to_topic(producer, email, "errors", "error_data_deduplicate");
					err != nil {
				logger.Error("failure in producing email to errors", zap.String("error", err.Error()))
			}
		}
		return
	}
	defer db.Close()

	query := `SELECT query.email_id, query.email_address, emails.email_id IS NULL as is_new
			  FROM (VALUES `
	for i, email := range emails {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', '%s')", email["email_id"].(string), email["address"].(string))
	}
	query += `) AS query(email_id, email_address)
			  LEFT JOIN emails ON query.email_id = emails.email_id 
			  	AND query.email_address = emails.email_address;`
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("failure in querying DB", zap.String("error", err.Error()))
		for _, email := range emails {
			if err := produce_email_to_topic(producer, email, "errors", "error_data_deduplicate");
					err != nil {
				logger.Error("failure in producing email to errors", zap.String("error", err.Error()))
			}
		}
		return
	}
	defer rows.Close()
	for i := 0; i < len(emails); i++ {
		rows.Next()
		values := make([]interface{}, 3)
		if err := rows.Scan(&values[0], &values[1], &values[2]); err != nil {
			logger.Error("failure in scanning row", zap.String("error", err.Error()))
			continue
		}
		if values[2].(bool) {
			logger.Info("new email", zap.String("email_id", values[0].(string)), zap.String("email_address", values[1].(string)))
			err := produce_email_to_topic(producer, emails[i], "new_emails", "new_email")
			if err != nil {
				logger.Error("failure in producing email to new_emails", zap.String("error", err.Error()))
				err = produce_email_to_topic(producer, emails[i], "errors", "error_data_deduplicate")
				if err != nil {
					logger.Error("failure in producing email to errors", zap.String("error", err.Error()))
				}
			}
			// store new emails to DB
		}
	}
}

func main() {
	logger.InitLogger("log", "data-async", 1, 7, "INFO")
	defer logger.LogErrorStackTrace()
	ctx := context.TODO()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:29092"})
	if err != nil {
		logger.Error(err.Error())
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092", 
		"group.id": "deduplicate-email",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error(err.Error())
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"emails"}, nil)

	// Manually assign partitions to the beginning offsets
	partitions, _ := consumer.Assignment()
	for _, partition := range partitions {
		consumer.Seek(kafka.TopicPartition{
			Topic:     partition.Topic,
			Partition: partition.Partition,
			Offset:    kafka.OffsetBeginning,
		}, 500)
	}

	logger.Info("data-async server started")

	// set up goroutine pool
	sem := semaphore.NewWeighted(MAX_WORKER_POOL_SIZE)

	// set up channel to receive batch of messages
	ch := make(chan []kafka.Message, MAX_BATCH_SIZE)

	go batch_message_delivery(consumer, ch)

	go func() {
		for {
			producer.Flush(500)
			time.Sleep(PRODUCER_FLUSH_INTERVAL)
		}
	}()

	// run consumer
	for {
		messages := <- ch
		if (len(messages) == 0) {
			continue
		}
		logger.Info("received batch of messages", zap.Int("messages num", len(messages)))
		emails := make([]map[string]interface{}, 0)
		for _, message := range messages {
			var email map[string]interface{}
			if err := json.Unmarshal(message.Value, &email); err != nil {
				logger.Error("failure in parsing received message", zap.String("error", err.Error()))
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key: []byte("error_data_deduplicate"),
					Value: message.Value,
				}, nil)
				continue
			}
			emails = append(emails, email)
		}
		if err := sem.Acquire(ctx, 1); err != nil {
			logger.Error("failure in acquiring semaphore", zap.String("error", err.Error()))
			logger.Warn("dropping batch of messages because failing to acquire semaphore")
			for _, message := range messages {
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key: []byte("error_data_deduplicate"),
					Value: message.Value,
				}, nil)
			}
			continue
		}
		go process_emails(emails, producer, sem)
	}
}