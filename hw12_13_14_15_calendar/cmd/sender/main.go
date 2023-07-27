package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/configs"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/sender"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/amqp"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	log := logger.GetDefaultLog()
	cfg := mustConfig(log)
	log = mustAppLog(log, cfg.Logger)

	var repo usecase.SenderRepo
	if cfg.Storage.Type == configs.StorageSQL {
		pg, err := postgres.NewPgByURL(log, cfg.Storage.DB.DSN())
		if err != nil {
			log.Fatal("main - postgres.NewPgByURL", err)
		}
		defer pg.Close()

		repo = sqlstorage.NewPgRepo(pg, log)
	} else {
		repo = memorystorage.NewMemoryStorage(log)
	}

	handler := getHandler(log, repo)
	senderMQ := getSender(log, cfg.MQ)
	defer func() {
		if err := senderMQ.Close(); err != nil {
			log.Error("error close sender", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer cancel()

	if err := senderMQ.Run(ctx, handler); err != nil {
		log.Error("stop sender", err)
	}
}

func mustConfig(l logger.AppLog) configs.SenderConfig {
	var cfg configs.SenderConfig
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

func getHandler(log logger.AppLog, repo usecase.SenderRepo) *sender.SubscriberHandler {
	uc := usecase.NewSenderUC(log, repo)
	return sender.NewSubscriberHandler(log, uc)
}

func getSender(log logger.AppLog, cfg configs.MQConf) *sender.Sender {
	return sender.NewSender(log, amqp.NewClientAMQP(log, cfg.DSN(), cfg.Queue))
}
