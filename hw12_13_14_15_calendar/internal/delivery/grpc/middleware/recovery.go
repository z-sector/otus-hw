package middleware

import (
	"context"
	"fmt"
	"runtime"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func Recovery(log logger.AppLog) grpc.UnaryServerInterceptor {
	opts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) (err error) {
			stack := make([]byte, 64<<10)
			stack = stack[:runtime.Stack(stack, false)]

			msg := fmt.Sprintf("panic triggered: %v", p)

			logWithRID := log.WithReqID(ctx)
			logWithRID.Base.Error().Any("error", string(stack)).Msg(msg)

			return status.Error(codes.Unknown, msg)
		}),
	}
	return recovery.UnaryServerInterceptor(opts...)
}
