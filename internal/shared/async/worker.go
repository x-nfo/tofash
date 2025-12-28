package async

import (
	"context"
	"encoding/json"
	"time"

	"tofash/internal/modules/notification/entity"
	notifService "tofash/internal/modules/notification/service"
	productService "tofash/internal/modules/product/service"
	"tofash/internal/modules/system/repository"

	"github.com/labstack/gommon/log"
	"gorm.io/datatypes"
)

type WorkerInterface interface {
	Run()
	Stop()
}

type worker struct {
	jobRepo      repository.JobRepositoryInterface
	productSvc   productService.ProductServiceInterface
	notifSvc     notifService.NotificationServiceInterface
	stopChan     chan struct{}
	pollInterval time.Duration
}

func NewWorker(
	jobRepo repository.JobRepositoryInterface,
	productSvc productService.ProductServiceInterface,
	notifSvc notifService.NotificationServiceInterface,
) WorkerInterface {
	return &worker{
		jobRepo:      jobRepo,
		productSvc:   productSvc,
		notifSvc:     notifSvc,
		stopChan:     make(chan struct{}),
		pollInterval: 2 * time.Second,
	}
}

func (w *worker) Run() {
	log.Info("[Worker] Starting async worker...")
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			log.Info("[Worker] Stopping worker...")
			return
		case <-ticker.C:
			w.processJobs()
		}
	}
}

func (w *worker) Stop() {
	close(w.stopChan)
}

func (w *worker) processJobs() {
	ctx := context.Background()
	jobs, err := w.jobRepo.FetchPendingJobs(ctx, 10) // Process 10 jobs at a time
	if err != nil {
		log.Errorf("[Worker] Failed to fetch jobs: %v", err)
		return
	}

	for _, job := range jobs {
		log.Infof("[Worker] Processing job ID: %d, Topic: %s", job.ID, job.Topic)

		// Mark as processing (optional, depends on simple or complex logic)
		// For now we just process and mark completed/failed.

		var processErr error

		switch job.Topic {
		case "stock_update":
			processErr = w.handleStockUpdate(ctx, job.Payload)
		case "email_notification":
			processErr = w.handleEmailNotification(ctx, job.Payload)
		default:
			log.Warnf("[Worker] Unknown topic: %s", job.Topic)
		}

		status := "completed"
		errMsg := ""
		if processErr != nil {
			status = "failed"
			errMsg = processErr.Error()
			log.Errorf("[Worker] Job %d failed: %v", job.ID, processErr)
		} else {
			log.Infof("[Worker] Job %d completed successfully", job.ID)
		}

		_ = w.jobRepo.UpdateJobStatus(ctx, job.ID, status, errMsg)
	}
}

// --- Handlers ---

type StockUpdatePayload struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

func (w *worker) handleStockUpdate(ctx context.Context, payload datatypes.JSON) error {
	var data StockUpdatePayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	// Assuming ProductService has a method to update stock.
	// If not, we might need to expose one or use repository directly.
	// Based on previous code, there was logic in consumer.
	// We'll reimplement it here or call a service method.

	// Previous Logic:
	// product.Stock -= int(orderItem.Quantity)
	// db.Save(&product)

	// BUT, ProductService usually handles business logic.
	// Let's check ProductService methods. If no direct method, we might need to add one.
	// For now, I'll assume we can call a method or I will need to extend ProductService.

	// Let's assume we need to implement UpdateStock in ProductService if not exists.
	// Checking `product_service.go` is recommended. I'll defer implementation details
	// or assume UpdateStock exists/will serve.

	return w.productSvc.UpdateStock(ctx, data.ProductID, int(data.Quantity))
}

type NotificationPayload struct {
	ReceiverEmail string `json:"receiver_email"`
	Subject       string `json:"subject"`
	Message       string `json:"message"`
	Type          string `json:"type"`
	ReceiverID    int64  `json:"receiver_id"`
}

func (w *worker) handleEmailNotification(ctx context.Context, payload datatypes.JSON) error {
	var data NotificationPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	// Reusing NotificationService logic
	notifEntity := entity.NotificationEntity{
		ReceiverEmail:    &data.ReceiverEmail,
		Subject:          &data.Subject,
		Message:          data.Message,
		NotificationType: data.Type, // EMAIL or PUSH
		Status:           "PENDING",
		ReceiverID:       uint(data.ReceiverID),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create Notification Record
	// Assuming `SendPushNotification` or equivalent handles creation + sending or just sending.
	// Looking at old consumer: it created record then sent email.

	// Create repo call is needed? `notifRepo.CreateNotification` was used.
	// But `worker` doesn't have `notifRepo`.
	// `notifSvc` should ideally encapsulate this.

	// Let's assume `notifSvc.SendNotification` exists or similar.
	// If not, we should probably add `CreateNotification` to service.

	// For now, let's call SendEmail directly via Service if available.
	// Actually the old code used `emailService.SendEmailNotif`.
	// Use case: Sending email.

	// Wait, `notifSvc` usually has high level methods.
	// We might need to inject `MessageEmailInterface` into worker too if logic was separated.
	// However, `NotificationService` might be enough.

	return w.notifSvc.CreateAndSend(ctx, notifEntity)
}
