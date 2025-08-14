package app

import (
	"context"
	"log"
	"notification/internal/kafka"
	"notification/internal/models"
	"notification/internal/service"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	emailService   *service.EmailService
	kafkaConsumer  *kafka.KafkaConsumer
	workerPoolSize int
}

func NewApp(emailService *service.EmailService, kafkaConsumer *kafka.KafkaConsumer, workerPoolSize int) *App {
	return &App{
		emailService:   emailService,
		kafkaConsumer:  kafkaConsumer,
		workerPoolSize: workerPoolSize,
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgChan := make(chan models.EmailMessage, 100)

	go func() {
		if err := a.kafkaConsumer.Consume(ctx, msgChan); err != nil {
			log.Fatalf("Ошибка в Kafka consumer: %v", err)
		}
	}()

	for i := 0; i < a.workerPoolSize; i++ {
		go a.emailWorker(ctx, msgChan, i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Получен сигнал завершения, останавливаем сервис...")
}

func (a *App) emailWorker(ctx context.Context, msgChan <-chan models.EmailMessage, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Воркер %d завершает работу", workerID)
			return
		case msg := <-msgChan:
			log.Printf("Воркер %d обрабатывает письмо для %v", workerID, msg.To)

			err := a.emailService.SendEmail(
				msg.To,
				msg.Subject,
				msg.TextBody,
				msg.HtmlBody,
				msg.Attachments,
			)

			if err != nil {
				log.Printf("Воркер %d: ошибка при отправке письма: %v", workerID, err)
			} else {
				log.Printf("Воркер %d: письмо успешно отправлено", workerID)
			}
		}
	}
}
