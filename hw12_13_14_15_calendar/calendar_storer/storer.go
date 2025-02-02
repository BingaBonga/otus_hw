package main

//nolint:depguard
import (
	"context"
	"errors"
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
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../calendar_storer/config.toml", "Path to configuration file")
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		consumer := kafka.Consumer{Config: config.Kafka, Logger: *logg, Context: ctx, Storage: storage, ConsumeMessage: consumeMessage}
		err := consumer.RegisterMessageConsumer()
		if err != nil {
			logg.Error("register message consumer failed", zap.Error(err))
			return
		}
	}()

	wg.Wait()
}

func consumeMessage(ctx context.Context, repository app.Storage, message *sarama.ConsumerMessage) error {
	var event storage.Event
	err := jsoniter.Unmarshal(message.Value, &event)
	if err != nil {
		return err
	}

	err = validateEvent(&event)
	if err != nil {
		return err
	}

	if event.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		event.ID = id.String()
	}

	event.StartDate = event.StartDate.UTC()
	err = repository.CreateEvent(ctx, &event)
	if err != nil {
		return err
	}

	return nil
}

func validateEvent(event *storage.Event) error {
	if event.Owner == "" {
		return errors.New("owner is required")
	}

	if len(event.Owner) > 256 {
		return errors.New("owner length can't be greater than 256")
	}

	if event.Title == "" {
		return errors.New("title is required")
	}

	if len(event.Title) > 256 {
		return errors.New("title length can't be greater than 256")
	}

	if event.Duration == 0 {
		return errors.New("duration is required")
	}

	if event.StartDate.Equal(time.Time{}) {
		return errors.New("startDate is required")
	}

	return nil
}
