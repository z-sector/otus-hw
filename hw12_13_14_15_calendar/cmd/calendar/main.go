package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/configs"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/grpc"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/http"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/hserver"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

var configFile string

type ServerI interface {
	Run()
	GetAddr() string
	Notify() <-chan error
	Shutdown() error
}

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	log := logger.GetDefaultLog()
	cfg := mustConfig(log)
	appLogg := mustAppLog(log, cfg.Logger)

	var eRepo usecase.EventRepo
	var iRepo usecase.HealthCheckRepo
	if cfg.Storage.Type == configs.StorageSQL {
		pg, err := postgres.NewPgByURL(appLogg, cfg.Storage.DB.DSN())
		if err != nil {
			appLogg.Fatal("main - mustUseCases - postgres.NewPgByURL", err)
		}
		defer pg.Close()

		pgStorage := sqlstorage.NewPgRepo(pg, appLogg)
		eRepo, iRepo = pgStorage, pgStorage
	} else {
		memStorage := memorystorage.NewMemoryStorage(appLogg)
		eRepo, iRepo = memStorage, memStorage
	}
	eUC, iUC := usecase.NewEventUC(appLogg, eRepo), usecase.NewInternalUC(appLogg, iRepo)

	server := mustServer(appLogg, cfg.Server, eUC, iUC)
	appLogg.Info("calendar is running...")
	server.Run()
	appLogg.Info(fmt.Sprintf("%s server listen and serve %s", cfg.Server.Type, server.GetAddr()))

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer cancel()

	select {
	case <-ctx.Done():
		appLogg.Info("interrupt signal received")
	case err := <-server.Notify():
		appLogg.Error("server error", err)
	}

	if err := server.Shutdown(); err != nil {
		appLogg.Error("failed to shutdown server", err)
	}
}

func mustConfig(l logger.AppLog) configs.CalendarConfig {
	var cfg configs.CalendarConfig
	if err := configs.ParseConfig(configFile, &cfg); err != nil {
		l.Fatal("incorrect config file", err)
	}
	l.Info(fmt.Sprintf("%v", cfg))
	return cfg
}

func mustAppLog(l logger.AppLog, cfg configs.LoggerConf) logger.AppLog {
	appLogg, err := logger.InitLog(cfg.Level, cfg.JSON)
	if err != nil {
		l.Fatal("incorrect init log", err)
	}
	return appLogg
}

func mustServer(l logger.AppLog, cfg configs.ServerConf, eUC *usecase.EventUC, iUC *usecase.InternalUC) ServerI {
	switch cfg.Type {
	case configs.ServerHTTP:
		handler := http.NewHandler(l, eUC, iUC)
		return hserver.NewHTTPServer(handler, hserver.Addr(cfg.Host, cfg.Port))
	case configs.ServerGRPC:
		return grpc.NewGrpcServer(cfg, l, eUC, iUC)
	}

	l.Base.Fatal().Msg("incorrect type server")
	return nil
}
