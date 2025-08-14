package main

import (
	"fmt"
	"notification/internal/app"
	"notification/internal/kafka"
	"notification/internal/service"
	"os"
	"strconv"
)

func main() {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	var kafkaAddr []string
	kafkaAddr = append(kafkaAddr, fmt.Sprintf("%s:%s", kafkaHost, kafkaPort))

	emailFrom := os.Getenv("EMAIL_FROM")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	emailService := service.NewEmailService(emailFrom, emailPassword, smtpHost, smtpPort)
	kafkaConsumer := kafka.NewKafkaConsumer(
		kafkaAddr,
		kafkaTopic,
	)

	app := app.NewApp(emailService, kafkaConsumer, 5)
	app.Run()
}
