package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ delivery.EventUCI = (*EventUC)(nil)

type EventRepo interface {
	Create(ctx context.Context, data dto.CreateEventDTO) (internal.Event, error)
	Update(ctx context.Context, e *internal.Event) error
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

func (u *EventUC) CreateEvent(ctx context.Context, data dto.CreateEventDTO) (internal.Event, error) {
	if err := data.Validate(); err != nil {
		return internal.Event{}, err
	}

	return u.repo.Create(ctx, data)
}

func (u *EventUC) UpdateEvent(ctx context.Context, data dto.UpdateEventDTO) (internal.Event, error) {
	if err := data.Validate(); err != nil {
		return internal.Event{}, err
	}

	event, err := u.repo.GetByID(ctx, data.ID)
	if err != nil {
		return internal.Event{}, err
	}

	if data.LastVersion != event.Version {
		return internal.Event{}, internal.ErrStorageConflict
	}

	data.UpdateEvent(&event)
	if err := u.repo.Update(ctx, &event); err != nil {
		return internal.Event{}, err
	}

	return event, nil
}

func (u *EventUC) DeleteEvent(ctx context.Context, ID uuid.UUID) error {
	return u.repo.Delete(ctx, ID)
}

func (u *EventUC) GetByID(ctx context.Context, ID uuid.UUID) (internal.Event, error) {
	return u.repo.GetByID(ctx, ID)
}

func (u *EventUC) GetByPeriod(ctx context.Context, from, to time.Time) ([]internal.Event, error) {
	return u.repo.GetByPeriod(ctx, from, to)
}
