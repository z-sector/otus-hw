package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
)

type CreateEventDTO struct {
	Title            string    `binding:"required"`
	BeginTime        time.Time `binding:"required"`
	EndTime          time.Time `binding:"required"`
	Description      string
	UserID           uuid.UUID `binding:"required"`
	NotificationTime *time.Time
}

func (c CreateEventDTO) Validate() error {
	if err := ValidateTitle(c.Title); err != nil {
		return err
	}
	if err := ValidateBeginEndTime(c.BeginTime, c.EndTime); err != nil {
		return err
	}
	return ValidateNotifTime(c.NotificationTime, c.BeginTime)
}

type UpdateEventDTO struct {
	ID               uuid.UUID `json:"-"`
	Title            string    `binding:"required"`
	BeginTime        time.Time `binding:"required"`
	EndTime          time.Time `binding:"required"`
	Description      string
	UserID           uuid.UUID `binding:"required"`
	NotificationTime *time.Time
	LastVersion      int32 `binding:"required"`
}

func (u UpdateEventDTO) Validate() error {
	if err := ValidateTitle(u.Title); err != nil {
		return err
	}
	if err := ValidateBeginEndTime(u.BeginTime, u.EndTime); err != nil {
		return err
	}
	return ValidateNotifTime(u.NotificationTime, u.BeginTime)
}

func (u UpdateEventDTO) UpdateEvent(event *internal.Event) {
	event.UserID = u.UserID
	event.Title = u.Title
	event.BeginTime = u.BeginTime
	event.EndTime = u.EndTime
	event.NotificationTime = u.NotificationTime
	event.Description = u.Description
	event.NotifyStatus = 0
}
