package delivery

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
)

type EventUCI interface {
	CreateEvent(ctx context.Context, data dto.CreateEventDTO) (internal.Event, error)
	UpdateEvent(ctx context.Context, data dto.UpdateEventDTO) (internal.Event, error)
	DeleteEvent(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (internal.Event, error)
	GetByPeriod(ctx context.Context, from, to time.Time) ([]internal.Event, error)
}

type HealthCheckUCI interface {
	Ping(ctx context.Context) error
}
