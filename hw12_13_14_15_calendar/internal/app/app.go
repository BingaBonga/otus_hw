package app

import (
	"context"
	"errors"
	"time"

	//nolint:depguard
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/google/uuid"
)

type App struct {
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventsByPeriod(ctx context.Context, owner string, startTime time.Time, endTime time.Time) ([]storage.Event, error)
	GetEvents(ctx context.Context) ([]storage.Event, error)
}

func New(storage Storage) *App {
	return &App{storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	err := a.validateEvent(event)
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

	event.IsSend = false
	event.StartDate = event.StartDate.UTC()
	err = a.storage.CreateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	if event.ID == "" {
		return errors.New("id is required")
	}

	err := a.validateEvent(event)
	if err != nil {
		return err
	}

	event.IsSend = false
	event.StartDate = event.StartDate.UTC()
	err = a.storage.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventsDay(ctx context.Context, owner string) ([]storage.Event, error) {
	timeStart := time.Now().UTC().Truncate(24 * time.Hour)
	return a.storage.GetEventsByPeriod(ctx, owner, timeStart, timeStart.Add(24*time.Hour))
}

func (a *App) GetEventsWeek(ctx context.Context, owner string) ([]storage.Event, error) {
	timeStart := time.Now().UTC().Truncate(7 * 24 * time.Hour)
	return a.storage.GetEventsByPeriod(ctx, owner, timeStart, timeStart.Add(7*24*time.Hour))
}

func (a *App) GetEventsMonth(ctx context.Context, owner string) ([]storage.Event, error) {
	timeNow := time.Now().UTC()
	timeStart := time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.UTC)
	timeEnd := time.Date(timeNow.Year(), timeNow.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	return a.storage.GetEventsByPeriod(ctx, owner, timeStart, timeEnd)
}

func (a *App) validateEvent(event *storage.Event) error {
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
