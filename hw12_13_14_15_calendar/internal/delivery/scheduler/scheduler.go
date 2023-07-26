package scheduler

import (
	"context"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type SchedUCI interface {
	DeleteOldEvents(ctx context.Context, days int) error
	GetEventNotifications(ctx context.Context) ([]internal.EventNotification, error)
	SetProcessingNotifyStatus(ctx context.Context, notif internal.EventNotification) error
	SetNotSentStatus(ctx context.Context, notif internal.EventNotification) error
}

type Producer interface {
	Publish(ctx context.Context, notif internal.EventNotification) error
}

type SchedHandler struct {
	log      logger.AppLog
	uc       SchedUCI
	producer Producer
}

func NewSchedHandler(log logger.AppLog, uc SchedUCI, producer Producer) *SchedHandler {
	return &SchedHandler{log: log, uc: uc, producer: producer}
}

func (h *SchedHandler) DeleteOldEvents(days int) {
	h.log.Info("START delete old events")
	defer func() {
		h.log.Info("END delete old events")
	}()

	if err := h.uc.DeleteOldEvents(context.Background(), days); err != nil {
		h.log.Error("error delete old events", err)
	}
}

func (h *SchedHandler) SendNotification() {
	h.log.Info("START send notification")
	defer func() {
		h.log.Info("END send notification")
	}()

	ctx := context.Background()
	notifications, err := h.uc.GetEventNotifications(ctx)
	if err != nil {
		h.log.Error("error send notification", err)
		return
	}

	for _, n := range notifications {
		if err := h.uc.SetProcessingNotifyStatus(ctx, n); err != nil {
			h.log.Error("error set status", err)
			return
		}
		err := h.producer.Publish(ctx, n)
		if err != nil {
			h.log.Error("error send notification", err)
			if err := h.uc.SetNotSentStatus(ctx, n); err != nil {
				h.log.Error("error set status", err)
			}
			return
		}
	}
}
