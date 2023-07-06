package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/pb"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func TestUpdateEvent(t *testing.T) {
	ctx := context.Background()
	log := logger.Nop()

	memStorage := memorystorage.NewMemoryStorage(log)
	eUC := usecase.NewEventUC(log, memStorage)
	handler := NewEventGrpcHandler(log, eUC)

	t.Run("success", func(t *testing.T) {
		beginTime := time.Now().UTC()
		endTime := beginTime.Add(time.Second)
		notifTime := beginTime.Add(-time.Second)

		req := pb.CreateEventReq{
			Title:            "title",
			BeginTime:        timestamppb.New(beginTime),
			EndTime:          timestamppb.New(endTime),
			Description:      "description",
			UserId:           uuid.New().String(),
			NotificationTime: timestamppb.New(notifTime),
		}

		event, err := handler.CreateEvent(ctx, &req)
		require.NoError(t, err)

		newTitle := req.Title + "new"
		newEvent, err := handler.UpdateEvent(ctx, &pb.UpdateEventReq{
			Id:               event.Id,
			Title:            newTitle,
			BeginTime:        event.BeginTime,
			EndTime:          event.EndTime,
			Description:      event.Description,
			UserId:           event.UserId,
			NotificationTime: event.NotificationTime,
			LastVersion:      event.Version,
		})
		require.NoError(t, err)
		require.Equal(t, newTitle, newEvent.Title)
	})

	t.Run("error", func(t *testing.T) {
		beginTime := time.Now().UTC()
		endTime := beginTime.Add(time.Second)
		notifTime := beginTime.Add(-time.Second)

		req := pb.CreateEventReq{
			Title:            "title",
			BeginTime:        timestamppb.New(beginTime),
			EndTime:          timestamppb.New(endTime),
			Description:      "description",
			UserId:           uuid.New().String(),
			NotificationTime: timestamppb.New(notifTime),
		}

		event, err := handler.CreateEvent(ctx, &req)
		require.NoError(t, err)

		testcases := []struct {
			req     *pb.UpdateEventReq
			rpcCode codes.Code
		}{
			{
				req: &pb.UpdateEventReq{
					Id:               "",
					Title:            event.Title,
					BeginTime:        event.BeginTime,
					EndTime:          event.EndTime,
					Description:      event.Description,
					UserId:           event.UserId,
					NotificationTime: event.NotificationTime,
					LastVersion:      event.Version,
				},
				rpcCode: codes.InvalidArgument,
			},
			{
				req: &pb.UpdateEventReq{
					Id:               event.Id,
					Title:            event.Title,
					BeginTime:        event.BeginTime,
					EndTime:          event.EndTime,
					Description:      event.Description,
					UserId:           "",
					NotificationTime: event.NotificationTime,
					LastVersion:      event.Version,
				},
				rpcCode: codes.InvalidArgument,
			},
			{
				req: &pb.UpdateEventReq{
					Id:               event.Id,
					Title:            event.Title,
					BeginTime:        nil,
					EndTime:          event.EndTime,
					Description:      event.Description,
					UserId:           event.UserId,
					NotificationTime: event.NotificationTime,
					LastVersion:      event.Version,
				},
				rpcCode: codes.InvalidArgument,
			},
			{
				req: &pb.UpdateEventReq{
					Id:               event.Id,
					Title:            event.Title,
					BeginTime:        event.BeginTime,
					EndTime:          nil,
					Description:      event.Description,
					UserId:           event.UserId,
					NotificationTime: event.NotificationTime,
					LastVersion:      event.Version,
				},
				rpcCode: codes.InvalidArgument,
			},
			{
				req: &pb.UpdateEventReq{
					Id:               event.Id,
					Title:            event.Title,
					BeginTime:        event.BeginTime,
					EndTime:          event.EndTime,
					Description:      event.Description,
					UserId:           event.UserId,
					NotificationTime: event.NotificationTime,
					LastVersion:      event.Version + 1,
				},
				rpcCode: codes.Aborted,
			},
		}

		for i, tc := range testcases {
			tc := tc
			t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
				_, err := handler.UpdateEvent(ctx, tc.req)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.rpcCode, st.Code())
			})
		}
	})
}
