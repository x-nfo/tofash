package repository

import (
	"context"
	"errors"
	"time"

	"tofash/internal/modules/system/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobRepositoryInterface interface {
	CreateJob(ctx context.Context, topic string, payload interface{}) error
	FetchPendingJobs(ctx context.Context, limit int) ([]model.Job, error)
	UpdateJobStatus(ctx context.Context, jobID uint, status string, errorMsg string) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepositoryInterface {
	return &jobRepository{db: db}
}

func (r *jobRepository) CreateJob(ctx context.Context, topic string, payload interface{}) error {
	jsonPayload, ok := payload.(datatypes.JSON)
	if !ok {
		// Try to cast to byte array if valid JSON byte
		if bytes, isBytes := payload.([]byte); isBytes {
			jsonPayload = datatypes.JSON(bytes)
		} else {
			// Fallback: This simple repo expects datatypes.JSON or compatible.
			// In production you might want json.Marshal here if payload is struct.
			// But for now let's assume valid datatypes.JSON is passed or let it fail.
			return errors.New("invalid payload type, expected datatypes.JSON or []byte")
		}
	}

	job := model.Job{
		Topic:   topic,
		Payload: jsonPayload,
		Status:  "pending",
	}
	return r.db.WithContext(ctx).Create(&job).Error
}

func (r *jobRepository) FetchPendingJobs(ctx context.Context, limit int) ([]model.Job, error) {
	var jobs []model.Job

	// Postgres FOR UPDATE SKIP LOCKED ensures concurrency safety
	err := r.db.WithContext(ctx).
		Clauses(gorm.Expr("FOR UPDATE SKIP LOCKED")).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&jobs).Error

	return jobs, err
}

func (r *jobRepository) UpdateJobStatus(ctx context.Context, jobID uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	return r.db.WithContext(ctx).
		Model(&model.Job{}).
		Where("id = ?", jobID).
		Updates(updates).Error
}
