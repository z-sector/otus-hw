package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/scheduler"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ scheduler.SchedUCI = &SchedulerUC{}

type SchedulerRepo interface {
	DeleteOldEvents(ctx context.Context, to time.Time) error
	GetEventsForNotify(ctx context.Context, time time.Time) ([]internal.Event, error)
	SetProcessingNotifyStatus(ctx context.Context, ID uuid.UUID) error
	SetNotSentStatus(ctx context.Context, ID uuid.UUID) error
}

type SchedulerUC struct {
	log  logger.AppLog
	repo SchedulerRepo
}

func NewSchedulerUC(log logger.AppLog, repo SchedulerRepo) *SchedulerUC {
	return &SchedulerUC{log: log, repo: repo}
}

func (u *SchedulerUC) DeleteOldEvents(ctx context.Context, days int) error {
	to := time.Now().UTC().AddDate(0, 0, -days)
	return u.repo.DeleteOldEvents(ctx, to)
}

func (u *SchedulerUC) GetEventNotifications(ctx context.Context) ([]internal.EventNotification, error) {
	timeNow := time.Now().UTC()
	events, err := u.repo.GetEventsForNotify(ctx, timeNow)
	if err != nil {
		return nil, err
	}

	notifications := make([]internal.EventNotification, len(events))
	for i := range events {
		notifications[i] = events[i].CreateNotification()
	}

	return notifications, nil
}

func (u *SchedulerUC) SetProcessingNotifyStatus(ctx context.Context, notif internal.EventNotification) error {
	return u.repo.SetProcessingNotifyStatus(ctx, notif.EventID)
}

func (u *SchedulerUC) SetNotSentStatus(ctx context.Context, notif internal.EventNotification) error {
	return u.repo.SetNotSentStatus(ctx, notif.EventID)
}
