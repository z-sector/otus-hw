package internal

import (
	"time"

	"github.com/google/uuid"
)

type RemindStatus int32

const (
	NotSentStatus RemindStatus = iota
	ProcessingStatus
	SentStatus
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
	NotifyStatus     RemindStatus `json:"-"`
}

func (e Event) CreateNotification() EventNotification {
	return EventNotification{
		EventID:    e.ID,
		EventTitle: e.Title,
		BeginTime:  e.BeginTime,
		UserID:     e.UserID,
	}
}

type EventNotification struct {
	EventID    uuid.UUID
	EventTitle string
	BeginTime  time.Time
	UserID     uuid.UUID
}
