package sender

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	mock_sender "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/sender/mocks"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func TestSubscriberHandler_Consume(t *testing.T) {
	mockUC := mock_sender.NewMockUCSenderI(gomock.NewController(t))
	handler := NewSubscriberHandler(logger.Nop(), mockUC)

	err := handler.Consume(context.Background(), []byte("invalid"))
	require.Error(t, err)

	notif := internal.EventNotification{
		EventID:    uuid.New(),
		EventTitle: "test",
		BeginTime:  time.Now().UTC(),
		UserID:     uuid.New(),
	}
	data, err := json.Marshal(notif)
	require.NoError(t, err)
	mockUC.EXPECT().SetSentNotifyStatus(gomock.Any(), gomock.Eq(notif)).Return(nil).Times(1)
	err = handler.Consume(context.Background(), data)
	require.NoError(t, err)
}
