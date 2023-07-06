package usecase

import (
	"context"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ delivery.HealthCheckUCI = (*InternalUC)(nil)

type HealthCheckRepo interface {
	Ping(ctx context.Context) error
}

type InternalUC struct {
	log  logger.AppLog
	repo HealthCheckRepo
}

func NewInternalUC(log logger.AppLog, repo HealthCheckRepo) *InternalUC {
	return &InternalUC{log: log, repo: repo}
}

func (i *InternalUC) Ping(ctx context.Context) error {
	return i.repo.Ping(ctx)
}
