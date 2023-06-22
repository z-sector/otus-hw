package middleware

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/requestid"
)

func customFields(ctx context.Context) logging.Fields {
	rid := requestid.CtxGetRequestID(ctx)
	return logging.Fields{logger.LogRIDKey, rid}
}

func Logging(log logger.AppLog) grpc.UnaryServerInterceptor {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadSent, logging.PayloadReceived),
		logging.WithFieldsFromContext(customFields),
	}

	logFunc := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.Base.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})

	return logging.UnaryServerInterceptor(logFunc, opts...)
}
