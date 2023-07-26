package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/configs"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/scheduler"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/amqp"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

const defaultShutdownTimeout = 20 * time.Second

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	log := logger.GetDefaultLog()
	cfg := mustConfig(log)
	log = mustAppLog(log, cfg.Logger)

	var repo usecase.SchedulerRepo
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

	clientAMQP := amqp.NewClientAMQP(log, cfg.MQ.DSN(), cfg.MQ.Queue)
	defer func() {
		if err := clientAMQP.Close(); err != nil {
			log.Error("error close client AMQP", err)
		}
	}()

	handler := getHandler(log, repo, clientAMQP)

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(cfg.Schedule.DeleteEventsCronbeat).Do(
		handler.DeleteOldEvents, cfg.Schedule.DeleteEventsAgeDay,
	)
	if err != nil {
		log.Fatal("failed cron task ", err)
	}
	if _, err := s.Every(cfg.Schedule.SendNotifyCronbeat).Do(handler.SendNotification); err != nil {
		log.Fatal("failed cron task ", err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer cancel()

	s.StartAsync()
	log.Info("scheduler is running...")

	<-ctx.Done()
	log.Info("interrupt signal received")

	if err := gracefulStop(s, defaultShutdownTimeout); err != nil {
		log.Error("failed to shutdown scheduler", err)
	}
}

func mustConfig(l logger.AppLog) configs.SchedulerConfig {
	var cfg configs.SchedulerConfig
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

func getHandler(log logger.AppLog, repo usecase.SchedulerRepo, amqp *amqp.ClientAMQP) *scheduler.SchedHandler {
	uc := usecase.NewSchedulerUC(log, repo)
	producer := scheduler.NewProducerRMQ(amqp)
	return scheduler.NewSchedHandler(log, uc, producer)
}

func gracefulStop(scheduler *gocron.Scheduler, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ok := make(chan struct{})
	go func() {
		scheduler.Stop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
