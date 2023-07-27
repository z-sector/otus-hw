package scheduler

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	mock_scheduler "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/scheduler/mocks"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func TestSchedHandler_DeleteOldEvents(t *testing.T) {
	mockUC, mockProducer := getDep(t)
	handler := NewSchedHandler(logger.Nop(), mockUC, mockProducer)
	days := 7

	mockUC.EXPECT().DeleteOldEvents(gomock.Any(), days).Times(1)
	handler.DeleteOldEvents(days)
}

func TestSchedHandler_SendNotification(t *testing.T) {
	mockUC, mockProducer := getDep(t)
	handler := NewSchedHandler(logger.Nop(), mockUC, mockProducer)

	e0 := createNotif("0")
	e1 := createNotif("1")

	gomock.InOrder(
		mockUC.EXPECT().GetEventNotifications(gomock.Any()).Return([]internal.EventNotification{e0, e1}, nil),
		mockUC.EXPECT().SetProcessingNotifyStatus(gomock.Any(), gomock.Eq(e0)).Return(nil),
		mockProducer.EXPECT().Publish(gomock.Any(), gomock.Eq(e0)).Return(nil),
		mockUC.EXPECT().SetProcessingNotifyStatus(gomock.Any(), gomock.Eq(e1)).Return(nil),
		mockProducer.EXPECT().Publish(gomock.Any(), gomock.Eq(e1)).Return(nil),
	)

	handler.SendNotification()
}

func getDep(t *testing.T) (*mock_scheduler.MockSchedUCI, *mock_scheduler.MockProducer) {
	t.Helper()

	ctrl := gomock.NewController(t)
	uc := mock_scheduler.NewMockSchedUCI(ctrl)
	producer := mock_scheduler.NewMockProducer(ctrl)
	return uc, producer
}

func createNotif(title string) internal.EventNotification {
	return internal.EventNotification{
		EventID:    uuid.New(),
		EventTitle: title,
		BeginTime:  time.Now().UTC(),
		UserID:     uuid.New(),
	}
}
