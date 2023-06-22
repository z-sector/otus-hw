package internal

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID
	Title            string
	BeginTime        time.Time
	EndTime          time.Time
	Description      string
	UserID           uuid.UUID
	NotificationTime *time.Time
	Version          int32
}
