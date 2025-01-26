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
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := configs.ReadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var calendar *app.App
	if config.DB.InMemory {
		storage := memorystorage.New()
		calendar = app.New(storage)
	} else {
		storage := sqlstorage.New()
		defer storage.Close(ctx)

		if err := storage.Connect(ctx, &config.DB); err != nil {
			logg.Error("database connection failed", zap.Error(err))
			return
		}

		logg.Info("database connection created")
		calendar = app.New(storage)
	}

	server := internalhttp.NewServer(ctx, logg, calendar)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("http server is stopping...")
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err := server.Start(config.HTTP); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	wg.Wait()
}
