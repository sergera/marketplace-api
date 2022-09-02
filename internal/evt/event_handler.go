package evt

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sergera/marketplace-api/internal/conf"
)

type topic uint8

const (
	Unconfirmed topic = iota
	InProgress
	Ready
)

var Topics map[topic]string = map[topic]string{
	Unconfirmed: "orders__unconfirmed",
	InProgress:  "orders__in_progress",
	Ready:       "orders__ready",
}

type EventHandler struct {
	producer *kafka.Producer
}

func NewEventHandler() *EventHandler {
	conf := conf.GetConf()
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": conf.KafkaHost + ":" + conf.KafkaPort})
	if err != nil {
		log.Panic(err)
	}
	eventHandler := EventHandler{p}
	go eventHandler.reportDeliveries()

	return &eventHandler
}

func (e *EventHandler) Close() {
	e.producer.Close()
}

func (e *EventHandler) Flush(timeoutMs int) {
	// Wait for message deliveries before shutting down
	e.producer.Flush(timeoutMs)
}

func (e *EventHandler) reportDeliveries() {
	// Delivery report handler for produced messages
	for e := range e.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (e *EventHandler) Produce(topic string, key string, value []byte) {
	e.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	},
		nil)
}
