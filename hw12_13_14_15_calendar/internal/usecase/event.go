package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/http"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ http.EventUCI = (*EventUC)(nil)

type EventRepo interface {
	Create(ctx context.Context, e *internal.Event) error
	Update(ctx context.Context, e internal.Event) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (internal.Event, error)
	GetByPeriod(ctx context.Context, from, to time.Time) ([]internal.Event, error)
}

type EventUC struct {
	log  logger.AppLog
	repo EventRepo
}

func NewEventUC(log logger.AppLog, repo EventRepo) *EventUC {
	return &EventUC{log: log, repo: repo}
}

func (u *EventUC) CreateEvent(ctx context.Context, ID, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
