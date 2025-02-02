package memorystorage

import (
	"context"
	"sync"
	"time"

	//nolint:depguard
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/storage"
)

type Storage struct {
	mu    sync.RWMutex
	event map[string]*storage.Event
}

func New() *Storage {
	return &Storage{event: make(map[string]*storage.Event)}
}

func (s *Storage) CreateEvent(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.event[event.ID] != nil {
		return storage.ErrEventAlreadyExist
	}

	s.event[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.event[event.ID] == nil {
		return storage.ErrEventDoesNotExist
	}

	s.event[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.event[id] == nil {
		return storage.ErrEventDoesNotExist
	}

	delete(s.event, id)
	return nil
}

func (s *Storage) GetEventsByPeriod(_ context.Context,
	owner string,
	startTime time.Time,
	endTime time.Time,
) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0)
	for _, e := range s.event {
		if e.Owner == owner && (e.StartDate.After(startTime) || e.StartDate.Equal(startTime)) && e.StartDate.Before(endTime) {
			events = append(events, *e)
		}
	}

	return events, nil
}

func (s *Storage) GetEvents(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, len(s.event))
	for _, e := range s.event {
		//nolint:makezero
		events = append(events, *e)
	}

	return events, nil
}
