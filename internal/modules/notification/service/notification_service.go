package service

import (
	"context"
	"tofash/internal/modules/notification/entity"
	"tofash/internal/modules/notification/repository"
)

type NotificationServiceInterface interface {
	GetAll(ctx context.Context, query entity.NotifyQueryString) ([]entity.NotificationEntity, int64, int64, error)
	GetByID(ctx context.Context, notifID uint) (*entity.NotificationEntity, error)
	MarkAsRead(ctx context.Context, notifID uint) error
	SendPushNotification(ctx context.Context, notification entity.NotificationEntity) error
}

type notificationService struct {
	repo repository.NotificationRepositoryInterface
}

func NewNotificationService(repo repository.NotificationRepositoryInterface) NotificationServiceInterface {
	return &notificationService{repo: repo}
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
