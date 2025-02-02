package main

//nolint:depguard
import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/kafka"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../calendar_scheduler/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := configs.ReadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage
	if config.DB.InMemory {
		storage = memorystorage.New()
	} else {
		storageSQL := sqlstorage.New()
		defer storageSQL.Close(ctx)

		if err := storageSQL.Connect(ctx, &config.DB); err != nil {
			logg.Error("database connection failed", zap.Error(err))
			return
		}

		logg.Info("database connection created")
		storage = storageSQL
	}

	producer, err := kafka.NewProducer(config.Kafka, logg)
	if err != nil {
		logg.Error("kafka producer creation failed", zap.Error(err))
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := startJob(ctx, storage, logg, producer, config.Schedule)
		if err != nil {
			logg.Error("cron job creation failed", zap.Error(err))
		}
	}()
	wg.Wait()
}

//nolint:lll
func startJob(ctx context.Context, storage app.Storage, logger *zap.Logger, producer *kafka.Producer, config configs.ScheduleConfig) error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	defer scheduler.Shutdown()
	_, err = scheduler.NewJob(gocron.CronJob(config.Cron, false), gocron.NewTask(clearEvents, ctx, storage, logger))
	if err != nil {
		return err
	}

	_, err = scheduler.NewJob(gocron.CronJob(config.Cron, false), gocron.NewTask(sendEvents, ctx, storage, logger, producer))
	if err != nil {
		return err
	}

	scheduler.Start()
	<-ctx.Done()
	return nil
}

func clearEvents(ctx context.Context, storage app.Storage, logger *zap.Logger) {
	logger.Info("clear event job: start")

	yearAgo := time.Now().AddDate(-1, 0, 0)
	events, err := storage.GetEvents(ctx)
	if err != nil {
		return
	}

	for _, event := range events {
		if event.StartDate.Before(yearAgo) {
			logger.Info("clear event job: delete event", zap.Any("event", event))

			err := storage.DeleteEvent(ctx, event.ID)
			if err != nil {
				logger.Error("clear event job: failed to delete event", zap.Error(err))
			}
		}
	}

	logger.Info("clear event job: end")
}

func sendEvents(ctx context.Context, storage app.Storage, logger *zap.Logger, producer *kafka.Producer) {
	logger.Info("send event job: start")

	timeNow := time.Now()
	events, err := storage.GetEvents(ctx)
	if err != nil {
		return
	}

	for _, event := range events {
		if !event.IsSend && event.StartDate.Add(time.Minute*time.Duration(event.RemindAt)).Before(timeNow) {
			err := producer.SendEventMessage(&event)
			if err != nil {
				logger.Error("send event job: failed to send event", zap.Error(err))
				continue
			}

			event.IsSend = true
			err = storage.UpdateEvent(ctx, &event)
			if err != nil {
				logger.Error("send event job: failed to update event", zap.Error(err))
				continue
			}
		}
	}

	logger.Info("send event job: end")
}
