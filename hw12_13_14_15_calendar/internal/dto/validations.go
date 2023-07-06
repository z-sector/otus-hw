package dto

import (
	"time"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
)

func ValidateTitle(title string) error {
	if title == "" {
		return &internal.ValidationError{Message: "Event Title cannot be empty"}
	}

	if len(title) > 100 {
		return &internal.ValidationError{Message: "Event Title length exceeded"}
	}

	return nil
}

func ValidateNotifTime(notif *time.Time, begin time.Time) error {
	if notif != nil && notif.IsZero() {
		return &internal.ValidationError{Message: "NotificationTime cannot be zero"}
	}

	if notif != nil && notif.After(begin) {
		return &internal.ValidationError{Message: "BeginTime must be greater or equal NotificationTime"}
	}
	return nil
}

func ValidateBeginEndTime(begin time.Time, end time.Time) error {
	if begin.IsZero() {
		return &internal.ValidationError{Message: "BeginTime cannot be zero"}
	}

	if end.IsZero() {
		return &internal.ValidationError{Message: "EndTime cannot be zero"}
	}

	if !begin.Before(end) {
		return &internal.ValidationError{Message: "EndTime must be greater BeginTime"}
	}

	return nil
}
