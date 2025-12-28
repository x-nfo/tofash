package rabbitmq

import (
	"context"
	"encoding/json"
	"tofash/internal/config"
	"tofash/internal/modules/notification/entity"
	"tofash/internal/modules/notification/message"
	"tofash/internal/modules/notification/repository"
	"tofash/internal/modules/notification/service"

	"github.com/labstack/gommon/log"
)

type ConsumeRabbitMQInterface interface {
	ConsumeMessage(queueName string) error
}

type consumeRabbitMQ struct {
	emailService        message.MessageEmailInterface
	notifRepository     repository.NotificationRepositoryInterface
	notificationService service.NotificationServiceInterface
}

// ConsumeMessage implements ConsumeRabbitMQInterface.
// ConsumeMessage implements ConsumeRabbitMQInterface.
func (c *consumeRabbitMQ) ConsumeMessage(queueName string) error {
	cfg := config.NewConfig()
	conn, err := config.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Errorf("[ConsumeMessage-1] Failed to connect to RabbitMQ: %v", err)
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[ConsumeMessage-2] Failed to open a channel: %v", err)
		return err
	}

	defer ch.Close()
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ConsumeMessage-3] Failed to consume messages: %v", err)
		return err
	}

	for d := range msgs {
		var notificationEntity entity.NotificationEntity
		log.Infof("Received a message: %s", d.Body)
		if err = json.Unmarshal(d.Body, &notificationEntity); err != nil {
			log.Errorf("Failed to unmarshal JSON: %v", err)
			continue
		}

		notificationEntity.Status = "PENDING"
		if notificationEntity.NotificationType == "EMAIL" {
			notificationEntity.Status = "SENT"
		}

		err = c.notifRepository.CreateNotification(context.Background(), notificationEntity)
		if err != nil {
			log.Errorf("Failed to create notification: %v", err)
			continue
		}

		go c.SendNotification(notificationEntity)
	}

	return nil
}

func (c *consumeRabbitMQ) SendNotification(notificationEntity entity.NotificationEntity) {
	switch notificationEntity.NotificationType {
	case "EMAIL":
		err := c.emailService.SendEmailNotif(*notificationEntity.ReceiverEmail, *notificationEntity.Subject, notificationEntity.Message)
		if err != nil {
			log.Errorf("Failed to send email notification: %v", err)
		}
	case "PUSH":
		c.notificationService.SendPushNotification(context.Background(), notificationEntity)
	}
}

func NewConsumeRabbitMQ(emailService message.MessageEmailInterface, notifRepository repository.NotificationRepositoryInterface, notificationService service.NotificationServiceInterface) ConsumeRabbitMQInterface {
	return &consumeRabbitMQ{
		emailService:        emailService,
		notifRepository:     notifRepository,
		notificationService: notificationService,
	}
}
