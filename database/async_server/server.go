package main

import (
	"context"
	"email-wizard/data/logger"
	"email-wizard/data/utils"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var MAX_EMAIL_BATCH_SIZE int = 20
var MAX_EVENT_BATCH_SIZE int = 50
var MAX_BATCH_WAIT_TIME = time.Second
var MAX_EMAIL_WORKER_POOL_SIZE int64 = 2
var MAX_EVENT_WORKER_POOL_SIZE int64 = 2
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
		Key:            []byte(key),
		Value:          email_json,
	}, nil)
	return nil
}

func batch_message_delivery(consumer *kafka.Consumer, max_batch_size int, ch chan []kafka.Message) {
	t0 := time.Now()
	buffer := make([]kafka.Message, 0)
	for {
		if time.Since(t0) > MAX_BATCH_WAIT_TIME || len(buffer) >= max_batch_size {
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
		logger.Error("failure in connecting to DB",
			zap.String("error", err.Error()),
			zap.String("producer", "deduplicate-email"))
		for _, email := range emails {
			if err := produce_email_to_topic(producer, email, "errors", "error_data_deduplicate"); err != nil {
				logger.Error("failure in producing email to errors",
					zap.String("error", err.Error()),
					zap.String("producer", "deduplicate-email"))
			}
		}
		return
	}
	defer db.Close()

	query := `SELECT query.email_id, query.email_address, emails.email_id IS NULL as is_new, query.seqid
			  FROM (VALUES `
	for i, email := range emails {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', '%s', %v)", email["email_id"].(string), email["address"].(string), i)
	}
	query += `) AS query(email_id, email_address, seqid)
			  LEFT JOIN emails ON query.email_id = emails.email_id 
			  	AND query.email_address = emails.email_address
			  ORDER BY query.seqid ASC;`
	rows, err := db.Query(query)
	if err != nil {
		logger.Error("failure in querying DB",
			zap.String("error", err.Error()),
			zap.String("producer", "deduplicate-email"))
		for _, email := range emails {
			if err := produce_email_to_topic(producer, email, "errors", "error_data_deduplicate"); err != nil {
				logger.Error("failure in producing email to errors",
					zap.String("error", err.Error()),
					zap.String("producer", "deduplicate-email"))
			}
		}
		return
	}
	defer rows.Close()
	for i := 0; i < len(emails); i++ {
		rows.Next()
		values := make([]interface{}, 4)
		if err := rows.Scan(&values[0], &values[1], &values[2], &values[3]); err != nil {
			logger.Error("failure in scanning row",
				zap.String("error", err.Error()),
				zap.String("producer", "deduplicate-email"))
			continue
		}
		if values[0].(string) != emails[i]["email_id"].(string) {
			logger.Error("email_id mismatch, should not happen!",
				zap.String("query", values[0].(string)),
				zap.String("email", emails[i]["email_id"].(string)),
				zap.String("producer", "deduplicate-email"))
			panic("email_id mismatch, should not happen!")
		}
		if values[2].(bool) {
			logger.Info("/new_email",
				zap.String("email_id", values[0].(string)),
				zap.String("email_address", values[1].(string)),
				zap.String("producer", "deduplicate-email"))
			// store new emails to DB
			_, err = utils.AddRow(map[string]interface{}{
				"user_id":          emails[i]["user_id"],
				"email_id":         emails[i]["email_id"],
				"email_address":    emails[i]["address"],
				"mailbox_type":     emails[i]["protocol"],
				"email_subject":    emails[i]["item"].(map[string]interface{})["subject"],
				"email_sender":     emails[i]["item"].(map[string]interface{})["sender"],
				"email_recipients": emails[i]["item"].(map[string]interface{})["recipient"],
				"email_date":       emails[i]["item"].(map[string]interface{})["date"],
				"email_content":    emails[i]["item"].(map[string]interface{})["content"],
				"event_ids":        []int32{},
			}, "emails")
			if err != nil {
				logger.Error("failure in storing new email to DB",
					zap.String("error", err.Error()),
					zap.String("email_id", values[0].(string)),
					zap.String("email_address", values[1].(string)),
					zap.String("producer", "deduplicate-email"))
				err = produce_email_to_topic(producer, emails[i], "errors", "error_data_deduplicate")
				if err != nil {
					logger.Error("failure in producing email to errors",
						zap.String("error", err.Error()),
						zap.String("producer", "deduplicate-email"))
				}
				continue
			}
			logger.Info("stored new email to DB",
				zap.String("email_id", values[0].(string)),
				zap.String("email_address", values[1].(string)),
				zap.String("producer", "deduplicate-email"))

			err := produce_email_to_topic(producer, emails[i], "new_emails", "new_email")
			if err != nil {
				logger.Error("failure in producing email to new_emails",
					zap.String("error", err.Error()),
					zap.String("email_id", values[0].(string)),
					zap.String("email_address", values[1].(string)),
					zap.String("producer", "deduplicate-email"))
				err = produce_email_to_topic(producer, emails[i], "errors", "error_data_deduplicate")
				if err != nil {
					logger.Error("failure in producing email to errors",
						zap.String("error", err.Error()),
						zap.String("producer", "deduplicate-email"))
				}
				continue
			}
			logger.Info("produced new email to new_emails",
				zap.String("email_id", values[0].(string)),
				zap.String("email_address", values[1].(string)),
				zap.String("producer", "deduplicate-email"))
		}
	}
	logger.Info("processed batch of emails", 
		zap.Int("emails num", len(emails)), 
		zap.String("producer", "deduplicate-email"))
}

func serve_deduplicate(producer *kafka.Producer, wg *sync.WaitGroup) {
	defer wg.Done()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "deduplicate-email",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("producer", "deduplicate-email"))
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"emails"}, nil)

	// // Manually assign partitions to the beginning offsets
	// partitions, _ := consumer.Assignment()
	// for _, partition := range partitions {
	// 	consumer.Seek(kafka.TopicPartition{
	// 		Topic:     partition.Topic,
	// 		Partition: partition.Partition,
	// 		Offset:    kafka.OffsetBeginning,
	// 	}, 500)
	// }

	logger.Info("data-async server started", zap.String("producer", "deduplicate-email"))

	// set up goroutine pool
	sem := semaphore.NewWeighted(MAX_EMAIL_WORKER_POOL_SIZE)

	// set up channel to receive batch of messages
	ch := make(chan []kafka.Message, MAX_EMAIL_BATCH_SIZE)

	go batch_message_delivery(consumer, MAX_EMAIL_BATCH_SIZE, ch)

	// run consumer
	ctx := context.TODO()
	for {
		messages := <-ch
		if len(messages) == 0 {
			continue
		}
		logger.Info("received batch of messages", zap.Int("messages num", len(messages)), zap.String("producer", "deduplicate-email"))
		emails := make([]map[string]interface{}, 0)
		for _, message := range messages {
			var email map[string]interface{}
			if err := json.Unmarshal(message.Value, &email); err != nil {
				logger.Error("failure in parsing received message", zap.String("error", err.Error()), zap.String("producer", "deduplicate-email"))
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key:            []byte("error_data_deduplicate"),
					Value:          message.Value,
				}, nil)
				continue
			}
			emails = append(emails, email)
		}
		if err := sem.Acquire(ctx, 1); err != nil {
			logger.Error("failure in acquiring semaphore", zap.String("error", err.Error()), zap.String("producer", "deduplicate-email"))
			logger.Warn("dropping batch of messages because failing to acquire semaphore", zap.String("producer", "deduplicate-email"))
			for _, message := range messages {
				err_topic := "errors"
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &err_topic, Partition: kafka.PartitionAny},
					Key:            []byte("error_data_deduplicate"),
					Value:          message.Value,
				}, nil)
			}
			continue
		}
		go process_emails(emails, producer, sem)
	}
}

func main() {
	logger.InitLogger("log", "data-async", 1, 7, "INFO")
	defer logger.LogErrorStackTrace()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:29092"})
	if err != nil {
		logger.Error(err.Error())
	}
	defer producer.Close()

	go func() {
		for {
			producer.Flush(500)
			time.Sleep(PRODUCER_FLUSH_INTERVAL)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(2)
	go serve_deduplicate(producer, &wg)
	go serve_store_events(producer, &wg)
	wg.Wait()
}
