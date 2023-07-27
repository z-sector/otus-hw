package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/sender"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ sender.UCSenderI = &SenderUC{}

type SenderRepo interface {
	SetSentNotifyStatus(ctx context.Context, ID uuid.UUID) error
}

type SenderUC struct {
	log  logger.AppLog
	repo SenderRepo
}

func NewSenderUC(log logger.AppLog, repo SenderRepo) *SenderUC {
	return &SenderUC{log: log, repo: repo}
}

func (u *SenderUC) SetSentNotifyStatus(ctx context.Context, notif internal.EventNotification) error {
	return u.repo.SetSentNotifyStatus(ctx, notif.EventID)
}
