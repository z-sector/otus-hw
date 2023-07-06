package grpc

import (
	"google.golang.org/grpc"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/configs"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/middleware"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc/pb"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/gserver"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func NewGrpcServer(
	cfg configs.ServerConf, l logger.AppLog, eUC *usecase.EventUC, iUC *usecase.InternalUC,
) *gserver.AppGRPCServer {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.RequestID, middleware.Logging(l), middleware.Recovery(l)),
	)

	server := gserver.NewGrpcServer(grpcServer, gserver.Addr(cfg.Host, cfg.Port))

	pb.RegisterInternalServiceServer(server.GRPC, NewInternalGrpcHandler(l, iUC))
	pb.RegisterEventServiceServer(server.GRPC, NewEventGrpcHandler(l, eUC))

	return server
}
