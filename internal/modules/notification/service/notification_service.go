package service

import (
	"context"
	"tofash/internal/modules/notification/entity"
	"tofash/internal/modules/notification/message"
	"tofash/internal/modules/notification/repository"

	"github.com/labstack/gommon/log"
)

type NotificationServiceInterface interface {
	GetAll(ctx context.Context, query entity.NotifyQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetByID(ctx context.Context, notifID uint) (*entity.NotificationEntity, error)
	MarkAsRead(ctx context.Context, notifID uint) error
	SendPushNotification(ctx context.Context, notification entity.NotificationEntity) error
	CreateAndSend(ctx context.Context, notification entity.NotificationEntity) error
}

type notificationService struct {
	repo     repository.NotificationRepositoryInterface
	emailSvc message.MessageEmailInterface
}

func NewNotificationService(repo repository.NotificationRepositoryInterface, emailSvc message.MessageEmailInterface) NotificationServiceInterface {
	return &notificationService{repo: repo, emailSvc: emailSvc}
}

func (s *notificationService) GetAll(ctx context.Context, query entity.NotifyQueryString) ([]entity.NotificationEntity, int64, int64, error) {
	return s.repo.GetAll(ctx, query)
}

func (s *notificationService) GetByID(ctx context.Context, notifID uint) (*entity.NotificationEntity, error) {
	return s.repo.GetByID(ctx, notifID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notifID uint) error {
	return s.repo.MarkAsRead(ctx, notifID)
}

func (s *notificationService) SendPushNotification(ctx context.Context, notification entity.NotificationEntity) error {
	// Logic to send push notification (e.g. Firebase)
	// For now, just placeholder
	// log.Println("Sending Push Notification:", notification)
	return nil
}

func (s *notificationService) CreateAndSend(ctx context.Context, notification entity.NotificationEntity) error {
	// 1. Create Notification Record (Pending)
	notification.Status = "PENDING"
	if notification.NotificationType == "EMAIL" {
		// Usually we mark as sent after sending, or pending before.
		// Old consumer logic: set status "SENT" if type EMAIL. Weird.
		// Let's stick to "PENDING" then update to "SENT".
	}

	err := s.repo.CreateNotification(ctx, notification)
	if err != nil {
		log.Errorf("[NotificationService] CreateNotification: %v", err)
		return err
	}

	// 2. Send Message
	switch notification.NotificationType {
	case "EMAIL":
		if notification.ReceiverEmail != nil && notification.Subject != nil {
			err = s.emailSvc.SendEmailNotif(*notification.ReceiverEmail, *notification.Subject, notification.Message)
			if err != nil {
				log.Errorf("[NotificationService] SendEmail: %v", err)
				// Don't fail the whole job? Or retry?
				// For now let's return error so job retries/fails.
				return err
			}
		}
	case "PUSH":
		// Reuse existing
		_ = s.SendPushNotification(ctx, notification)
	}

	return nil
}
