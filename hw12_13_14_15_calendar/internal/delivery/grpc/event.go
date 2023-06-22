package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/pb"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func NewEventGrpcHandler(log logger.AppLog, uc delivery.EventUCI) *EventGrpcHandler {
	return &EventGrpcHandler{log: log, uc: uc}
}

type EventGrpcHandler struct {
	log logger.AppLog
	uc  delivery.EventUCI
	pb.UnimplementedEventServiceServer
}

func (e *EventGrpcHandler) CreateEvent(ctx context.Context, req *pb.CreateEventReq) (*pb.Event, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil || userID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid for `user_id`")
	}

	if err := req.BeginTime.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid begin_time: %s", err.Error()))
	}
	if err := req.EndTime.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid end_time: %s", err.Error()))
	}
	if req.NotificationTime != nil {
		if err := req.NotificationTime.CheckValid(); err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid notification_time: %s", err.Error()))
		}
	}

	data := dto.CreateEventDTO{
		Title:            req.Title,
		BeginTime:        req.BeginTime.AsTime(),
		EndTime:          req.EndTime.AsTime(),
		Description:      req.Description,
		UserID:           userID,
		NotificationTime: tspToTime(req.NotificationTime),
	}

	event, err := e.uc.CreateEvent(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return eventToEventpb(event), nil
}

func (e *EventGrpcHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventReq) (*pb.Event, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil || ID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid for `id`")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil || userID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid for `user_id`")
	}

	if err := req.BeginTime.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid begin_time: %s", err.Error()))
	}
	if err := req.EndTime.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid end_time: %s", err.Error()))
	}
	if req.NotificationTime != nil {
		if err := req.NotificationTime.CheckValid(); err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid notification_time: %s", err.Error()))
		}
	}

	data := dto.UpdateEventDTO{
		ID:               ID,
		Title:            req.Title,
		BeginTime:        req.BeginTime.AsTime(),
		EndTime:          req.EndTime.AsTime(),
		Description:      req.Description,
		UserID:           userID,
		NotificationTime: tspToTime(req.NotificationTime),
		LastVersion:      req.LastVersion,
	}

	event, err := e.uc.UpdateEvent(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return eventToEventpb(event), nil
}

func (e *EventGrpcHandler) DeleteEvent(ctx context.Context, req *pb.EventIDReq) (*pb.DeleteEventResp, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil || ID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid for `id`")
	}

	if err := e.uc.DeleteEvent(ctx, ID); err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return &pb.DeleteEventResp{Ok: true}, nil
}

func (e *EventGrpcHandler) GetByID(ctx context.Context, req *pb.EventIDReq) (*pb.Event, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil || ID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid for `id`")
	}

	event, err := e.uc.GetByID(ctx, ID)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return eventToEventpb(event), nil
}

func (e *EventGrpcHandler) GetByPeriod(ctx context.Context, req *pb.EventPeriodReq) (*pb.EventList, error) {
	if err := req.From.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid from: %s", err.Error()))
	}
	if err := req.To.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid to: %s", err.Error()))
	}

	events, err := e.uc.GetByPeriod(ctx, req.GetFrom().AsTime(), req.GetTo().AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	res := make([]*pb.Event, len(events))
	for i, v := range events {
		res[i] = eventToEventpb(v)
	}

	return &pb.EventList{Items: res}, nil
}
