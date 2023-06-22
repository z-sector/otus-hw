package grpc

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/pb"
)

func tspToTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	res := t.AsTime()
	return &res
}

func timeToTsp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func eventToEventpb(event internal.Event) *pb.Event {
	return &pb.Event{
		Id:               event.ID.String(),
		Title:            event.Title,
		BeginTime:        timestamppb.New(event.BeginTime),
		EndTime:          timestamppb.New(event.EndTime),
		Description:      event.Description,
		UserId:           event.UserID.String(),
		NotificationTime: timeToTsp(event.NotificationTime),
		Version:          event.Version,
	}
}
