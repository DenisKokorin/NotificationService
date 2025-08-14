package kafka

import (
	"context"
	"encoding/json"
	"log"
	"notification/internal/models"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		Dialer:   dialer,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		reader: reader,
	}
}

func (kc *KafkaConsumer) Consume(ctx context.Context, msgChan chan<- models.EmailMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Ошибка при чтении сообщения: %v", err)
				continue
			}

			var emailMsg models.EmailMessage
			if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
				log.Printf("Ошибка при декодировании сообщения: %v", err)
				continue
			}

			msgChan <- emailMsg
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}
