package utils

import (
	"email-wizard/backend/logger"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var CONSUMER_WAIT_TIME_PER_MESSAGE = time.Millisecond * 500
var BUFFER_SIZE_PER_USER = 50

type Handler struct {
	rec     chan map[string]interface{}
	user_id int
}

type Hub struct {
	// register websocket handler
	handlers map[int]*Handler
	// register requests from websocket handler
	register chan *Handler
	// unregister requests from websocket handler
	unregister chan *Handler
	// kafka consumer channel
	message chan *kafka.Message
	// consumer
	consumer *kafka.Consumer
}

func (h *Hub) run() {
	defer h.consumer.Close()
	for {
		select {
		case handler := <-h.register:
			h.handlers[handler.user_id] = handler
			logger.Info("registered websocket handler", zap.Int("user_id", handler.user_id))
		case handler := <-h.unregister:
			if _, ok := h.handlers[handler.user_id]; ok {
				delete(h.handlers, handler.user_id)
				close(handler.rec)
			}
			logger.Info("unregistered websocket handler", zap.Int("user_id", handler.user_id))
		case msg := <-h.message:
			event := make(map[string]interface{})
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("failure in unmarshalling event",
					zap.String("error", err.Error()),
					zap.String("producer", "respond-events"))
				continue
			}
			user_id := event["user_id"].(int)
			if handler, ok := h.handlers[user_id]; ok {
				logger.Info("sending event to websocket handler", zap.Int("user_id", user_id))
				handler.rec <- event
			}
		}
	}
}

func (h *Hub) listen() {
	for {
		msg, err := h.consumer.ReadMessage(CONSUMER_WAIT_TIME_PER_MESSAGE)
		if err != nil && err.(kafka.Error).Code() != kafka.ErrTimedOut {
			logger.Error("failure in reading message from kafka",
				zap.String("error", err.Error()),
				zap.String("producer", "respond-events"))
			continue
		}
		if err != nil {
			continue
		}
		event := make(map[string]interface{})
		if err = json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("failure in unmarshalling event",
				zap.String("error", err.Error()),
				zap.String("producer", "respond-events"))
			continue
		}
		user_id := int(event["user_id"].(float64))
		if handler, ok := h.handlers[user_id]; ok {
			logger.Info("sending event to websocket handler", zap.Int("user_id", user_id))
			handler.rec <- event
		}
	}
}

func NewHub() *Hub {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "respond-events",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("producer", "respond-events"))
	}
	consumer.SubscribeTopics([]string{"events"}, nil)

	h := &Hub{
		handlers:   make(map[int]*Handler),
		register:   make(chan *Handler),
		unregister: make(chan *Handler),
		message:    make(chan *kafka.Message),
		consumer:   consumer,
	}
	go h.run()
	go h.listen()
	return h
}

func (h *Hub) RegisterHandler(user_id int) *Handler {
	handler := &Handler{
		rec:     make(chan map[string]interface{}, BUFFER_SIZE_PER_USER),
		user_id: user_id,
	}
	h.register <- handler
	return handler
}

func HandleWebSocketConnection(conn *websocket.Conn, user_id int, h *Hub) {
	handler := h.RegisterHandler(user_id)
	logger.Info("Websocket connection opened", zap.Int("user_id", user_id))
	// send stored events to clients
	events, err := GetUserEvents(user_id)
	if err != nil {
		logger.Error("failure in getting events from database", zap.String("error", err.Error()))
		close_msg := websocket.FormatCloseMessage(
			websocket.CloseInternalServerErr,
			"Internal server error: failure in getting events from database")
		conn.WriteMessage(websocket.CloseMessage, close_msg)
		return
	}
	go func() {
		for _, event := range events {
			handler.rec <- map[string]interface{}{
				"event": event,
			}
		}
	}()
	// send incrementing events to clients
	done_channel := make(chan bool)

	go func() {
		// Read message from browser
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Info("Websocket connection closed")
			} else {
				logger.Error("Websocket connection error",
					zap.String("error", err.Error()))
			}
			done_channel <- true
			return
		}

	}()
	for {
		select {
		case <-done_channel:
			h.unregister <- handler
			return
		case event := <-handler.rec:
			if err := conn.WriteJSON(event["event"]); err != nil {
				logger.Error("failure in writing message to websocket",
					zap.String("error", err.Error()))
				close_msg := websocket.FormatCloseMessage(
					websocket.CloseInternalServerErr,
					"Internal server error")
				conn.WriteMessage(websocket.CloseMessage, close_msg)
				return
			}
		}
	}
}
