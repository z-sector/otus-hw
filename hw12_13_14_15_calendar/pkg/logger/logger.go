package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/requestid"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

type AppLog struct {
	Base zerolog.Logger
}

func (l AppLog) Trace(msg string) {
	l.Base.Trace().Msg(msg)
}

func (l AppLog) Debug(msg string) {
	l.Base.Debug().Msg(msg)
}

func (l AppLog) Info(msg string) {
	l.Base.Info().Msg(msg)
}

func (l AppLog) Warn(msg string) {
	l.Base.Warn().Msg(msg)
}

func (l AppLog) Error(msg string, err error) {
	l.Base.Error().Err(err).Msg(msg)
}

func (l AppLog) Fatal(msg string, err error) {
	l.Base.Fatal().Err(err).Msg(msg)
}

func (l AppLog) Panic(msg string, err error) {
	l.Base.Panic().Err(err).Msg(msg)
}

func (l AppLog) WithReqID(ctx context.Context) AppLog {
	reqID := requestid.CtxGetRequestID(ctx)
	if reqID != "" {
		return AppLog{l.Base.With().Str("request_id", reqID).Logger()}
	}
	return l
}

func (l AppLog) IsDebugging() bool {
	return l.Base.GetLevel() <= zerolog.DebugLevel
}

func GetDefaultLog() AppLog {
	return AppLog{getLog(zerolog.DebugLevel, getConsoleWriter())}
}

func InitLog(levelStr string, json bool) (AppLog, error) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		return AppLog{}, err
	}

	var writer io.Writer
	if json {
		writer = os.Stdout
	} else {
		writer = getConsoleWriter()
	}

	return AppLog{getLog(level, writer)}, nil
}

func Nop() AppLog {
	return AppLog{zerolog.Nop()}
}

func getLog(level zerolog.Level, writer io.Writer) zerolog.Logger {
	return zerolog.New(writer).Level(level).With().Timestamp().Caller().Logger()
}

func getConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}
}
