package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ Subscriber = &SubscriberHandler{}

type UCSenderI interface {
	SetSentNotifyStatus(ctx context.Context, notif internal.EventNotification) error
}

type SubscriberHandler struct {
	log logger.AppLog
	uc  UCSenderI
}

func NewSubscriberHandler(log logger.AppLog, uc UCSenderI) *SubscriberHandler {
	return &SubscriberHandler{log: log, uc: uc}
}

func (h *SubscriberHandler) Consume(ctx context.Context, data []byte) error {
	var notif internal.EventNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return fmt.Errorf("SubscriberHandler - Consume - json.Unmarshal: %w", err)
	}
	h.log.Info(fmt.Sprintf("SEND EMAIL FOR EVENT WITH ID=%s", notif.EventID))
	if err := h.uc.SetSentNotifyStatus(ctx, notif); err != nil {
		return fmt.Errorf("SubscriberHandler - Consume - h.uc.SetSentNotifyStatus: %w", err)
	}
	return nil
}
