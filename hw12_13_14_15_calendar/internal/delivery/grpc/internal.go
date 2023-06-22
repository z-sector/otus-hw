package grpc

import (
	"context"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/pb"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func NewInternalGrpcHandler(log logger.AppLog, uc delivery.HealthCheckUCI) *InternalGrpcHandler {
	return &InternalGrpcHandler{log: log, uc: uc}
}

type InternalGrpcHandler struct {
	log logger.AppLog
	uc  delivery.HealthCheckUCI
	pb.UnimplementedInternalServiceServer
}

func (i *InternalGrpcHandler) Ping(ctx context.Context, _ *pb.PingReq) (*pb.PingResp, error) {
	res := pb.PingResp{}

	if err := i.uc.Ping(ctx); err == nil {
		res.Ok = true
	} else {
		res.Ok = false
		i.log.Error("failed ping", err)
	}

	return &res, nil
}
