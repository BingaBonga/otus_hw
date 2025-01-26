package memorystorage

import (
	"context"
	"testing"
	"time"

	//nolint:depguard
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	testEvent := &storage.Event{
		ID:          "test_id",
		Title:       "test_title",
		Owner:       "test_user",
		StartDate:   time.Now(),
		Duration:    30,
		Description: "test_description",
		RemindAt:    100,
	}

	t.Run("event created", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.CreateEvent(ctx, testEvent)

		require.NoError(t, err)
		events, err := memory.GetEventsByPeriod(ctx, testEvent.Owner, testEvent.StartDate, testEvent.StartDate.Add(1))
		require.NoError(t, err)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0].ID, "test_id")
		require.Equal(t, events[0].Title, "test_title")
		require.Equal(t, events[0].Owner, "test_user")
		require.Equal(t, events[0].Duration, time.Duration(30))
		require.Equal(t, events[0].Description, "test_description")
		require.Equal(t, events[0].RemindAt, int64(100))
		require.NotNil(t, events[0].StartDate)

		err = memory.DeleteEvent(ctx, events[0].ID)
		require.NoError(t, err)
	})

	t.Run("event create already exists", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.CreateEvent(ctx, testEvent)
		require.NoError(t, err)

		err = memory.CreateEvent(ctx, &storage.Event{
			ID: testEvent.ID,
		})
		require.Equal(t, err, storage.ErrEventAlreadyExist)

		err = memory.DeleteEvent(ctx, testEvent.ID)
		require.NoError(t, err)
	})

	t.Run("event updated", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.CreateEvent(ctx, testEvent)
		require.NoError(t, err)

		err = memory.UpdateEvent(ctx, &storage.Event{
			ID:          testEvent.ID,
			Title:       "test_title2",
			Owner:       "test_user2",
			StartDate:   time.Now(),
			Duration:    32,
			Description: "test_description2",
			RemindAt:    120,
		})
		require.NoError(t, err)

		events, err := memory.GetEventsByPeriod(ctx, "test_user2", testEvent.StartDate, testEvent.StartDate.Add(1))
		require.NoError(t, err)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0].ID, "test_id")
		require.Equal(t, events[0].Title, "test_title2")
		require.Equal(t, events[0].Owner, "test_user2")
		require.Equal(t, events[0].Duration, time.Duration(32))
		require.Equal(t, events[0].Description, "test_description2")
		require.Equal(t, events[0].RemindAt, int64(120))
		require.NotNil(t, events[0].StartDate)

		err = memory.DeleteEvent(ctx, events[0].ID)
		require.NoError(t, err)
	})

	t.Run("event update does not exists", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.UpdateEvent(ctx, &storage.Event{
			ID: "not_exists",
		})
		require.Equal(t, err, storage.ErrEventDoesNotExist)
	})

	t.Run("event deleted", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.CreateEvent(ctx, testEvent)
		require.NoError(t, err)

		err = memory.DeleteEvent(ctx, testEvent.ID)
		require.NoError(t, err)

		events, err := memory.GetEventsByPeriod(ctx, testEvent.Owner, testEvent.StartDate, testEvent.StartDate.Add(1))
		require.NoError(t, err)
		require.Equal(t, len(events), 0)
	})

	t.Run("event delete does not exists", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.DeleteEvent(ctx, "not_exists")
		require.Equal(t, err, storage.ErrEventDoesNotExist)
	})

	t.Run("events get by period", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memory := New()
		err := memory.CreateEvent(ctx, testEvent)
		require.NoError(t, err)

		events, err := memory.GetEventsByPeriod(ctx, testEvent.Owner, testEvent.StartDate, testEvent.StartDate.Add(1))
		require.NoError(t, err)
		require.Equal(t, len(events), 1)

		events, err = memory.GetEventsByPeriod(ctx, testEvent.Owner, testEvent.StartDate.Add(1), testEvent.StartDate.Add(2))
		require.NoError(t, err)
		require.Equal(t, len(events), 0)

		err = memory.DeleteEvent(ctx, testEvent.ID)
		require.NoError(t, err)
	})
}
