package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/configs"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/http"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/httpserver"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	log := logger.GetDefaultLog()

	cfg, err := configs.NewConfig(configFile)
	if err != nil {
		log.Fatal("incorrect config file", err)
	}
	log.Info(fmt.Sprintf("%v", cfg))

	appLogg, err := logger.InitLog(cfg.Logger.Level, cfg.Logger.JSON)
	if err != nil {
		log.Fatal("incorrect init log", err)
	}

	var eRepo usecase.EventRepo
	var iRepo usecase.HealthCheckRepo
	if cfg.Storage.Type == configs.StorageSQL {
		pg, err := postgres.NewPgByURL(appLogg, cfg.Storage.DB.DSN())
		if err != nil {
			appLogg.Fatal("app - Run - postgres.New", err)
		}
		defer pg.Close()

		pgStorage := sqlstorage.NewPgRepo(pg, appLogg)
		eRepo, iRepo = pgStorage, pgStorage
	} else {
		memStorage := memorystorage.NewMemoryStorage(appLogg)
		eRepo, iRepo = memStorage, memStorage
	}

	handler := http.NewHandler(
		appLogg,
		usecase.NewEventUC(appLogg, eRepo),
		usecase.NewInternalUC(appLogg, iRepo),
	)

	appLogg.Info("calendar is running...")
	server := httpserver.NewHTTPServer(handler, httpserver.Addr(cfg.HTTP.Host, cfg.HTTP.Port))
	server.Run()
	appLogg.Info(fmt.Sprintf("HTTP server listen and serve %s", server.GetAddr()))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	select {
	case <-ctx.Done():
		appLogg.Info("interrupt signal received")
	case err = <-server.Notify():
		appLogg.Error("server error", err)
	}

	err = server.Shutdown()
	if err != nil {
		appLogg.Error("failed to shutdown server", err)
	}
}
