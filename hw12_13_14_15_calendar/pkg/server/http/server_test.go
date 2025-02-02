package internalhttp

//nolint:depguard
import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/api"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/app"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	memorystorage "github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/pkg/storage/memory"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

//nolint:funlen
func TestStorage(t *testing.T) {
	var level zapcore.Level
	logg, _ := logger.New(level, os.TempDir()+"/test.log")

	testEvent := &storage.Event{
		ID:          "test_id",
		Title:       "test_title",
		Owner:       "test_user",
		StartDate:   time.Now(),
		Duration:    30,
		Description: "test_description",
		RemindAt:    100,
	}
	testEventMarshal, _ := json.Marshal(&testEvent)

	t.Run("Test start server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		go func() {
			err := server.Start(configs.HTTPConfig{})
			require.NoError(t, err)
		}()

		cancel()
	})

	t.Run("Create event", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		req := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		respBody, _ := io.ReadAll(resp.Body)
		respEvent := &storage.Event{}

		err := json.Unmarshal(respBody, respEvent)
		require.NoError(t, err)

		require.Equal(t, resp.Code, 200)
		require.Equal(t, testEvent.ID, respEvent.ID)
		require.Equal(t, testEvent.Title, respEvent.Title)
		require.Equal(t, testEvent.Owner, respEvent.Owner)
		require.Equal(t, testEvent.StartDate.UTC(), respEvent.StartDate)
		require.Equal(t, testEvent.Duration, respEvent.Duration)
		require.Equal(t, testEvent.Description, respEvent.Description)
		require.Equal(t, testEvent.RemindAt, respEvent.RemindAt)
	})

	t.Run("Update event", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		reqCreate := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		respCreate := httptest.NewRecorder()
		handler.ServeHTTP(respCreate, reqCreate)

		updatedTime := time.Now()
		updatedEvent := &storage.Event{
			ID:          testEvent.ID,
			Title:       "test_title2",
			Owner:       "test_user2",
			StartDate:   updatedTime,
			Duration:    32,
			Description: "test_description2",
			RemindAt:    120,
		}
		updatedEventMarshal, _ := json.Marshal(&updatedEvent)

		reqUpdate := httptest.NewRequest("PUT", "/event", bytes.NewBuffer(updatedEventMarshal))
		respUpdate := httptest.NewRecorder()
		handler.ServeHTTP(respUpdate, reqUpdate)

		respBody, _ := io.ReadAll(respUpdate.Body)
		respEvent := &storage.Event{}

		err := json.Unmarshal(respBody, respEvent)
		require.NoError(t, err)

		require.Equal(t, respUpdate.Code, 200)
		require.Equal(t, updatedEvent.ID, respEvent.ID)
		require.Equal(t, updatedEvent.Title, respEvent.Title)
		require.Equal(t, updatedEvent.Owner, respEvent.Owner)
		require.Equal(t, updatedEvent.StartDate.UTC(), respEvent.StartDate)
		require.Equal(t, updatedEvent.Duration, respEvent.Duration)
		require.Equal(t, updatedEvent.Description, respEvent.Description)
		require.Equal(t, updatedEvent.RemindAt, respEvent.RemindAt)
	})

	t.Run("Delete event", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		reqCreate := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		respCreate := httptest.NewRecorder()
		handler.ServeHTTP(respCreate, reqCreate)

		deleteEvent := &storage.Event{
			ID: testEvent.ID,
		}
		deleteEventMarshal, _ := json.Marshal(&deleteEvent)

		reqDelete := httptest.NewRequest("DELETE", "/event", bytes.NewBuffer(deleteEventMarshal))
		respDelete := httptest.NewRecorder()
		handler.ServeHTTP(respDelete, reqDelete)

		require.Equal(t, respDelete.Code, 200)
	})

	t.Run("Get day events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		reqCreate := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		respCreate := httptest.NewRecorder()
		handler.ServeHTTP(respCreate, reqCreate)

		reqGet := httptest.NewRequest("GET", "/event/test_user/getDay", bytes.NewBuffer(nil))
		respGet := httptest.NewRecorder()
		handler.ServeHTTP(respGet, reqGet)

		respBody, _ := io.ReadAll(respGet.Body)
		var respEvents []*storage.Event

		err := json.Unmarshal(respBody, &respEvents)
		require.NoError(t, err)

		require.Equal(t, respGet.Code, 200)
		require.Equal(t, len(respEvents), 1)
	})

	t.Run("Get week events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		reqCreate := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		respCreate := httptest.NewRecorder()
		handler.ServeHTTP(respCreate, reqCreate)

		reqGet := httptest.NewRequest("GET", "/event/test_user/getWeek", bytes.NewBuffer(nil))
		respGet := httptest.NewRecorder()
		handler.ServeHTTP(respGet, reqGet)

		respBody, _ := io.ReadAll(respGet.Body)
		var respEvents []*storage.Event

		err := json.Unmarshal(respBody, &respEvents)
		require.NoError(t, err)

		require.Equal(t, respGet.Code, 200)
		require.Equal(t, len(respEvents), 1)
	})

	t.Run("Get month events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calendar := app.New(memorystorage.New())
		server := NewServer(ctx, logg, calendar)
		require.NotNil(t, server)

		handler := api.HandlerFromMux(server, http.NewServeMux())
		require.NotNil(t, handler)

		ts := httptest.NewServer(handler)
		defer ts.Close()

		reqCreate := httptest.NewRequest("POST", "/event", bytes.NewBuffer(testEventMarshal))
		respCreate := httptest.NewRecorder()
		handler.ServeHTTP(respCreate, reqCreate)

		reqGet := httptest.NewRequest("GET", "/event/test_user/getMonth", bytes.NewBuffer(nil))
		respGet := httptest.NewRecorder()
		handler.ServeHTTP(respGet, reqGet)

		respBody, _ := io.ReadAll(respGet.Body)
		var respEvents []*storage.Event

		err := json.Unmarshal(respBody, &respEvents)
		require.NoError(t, err)

		require.Equal(t, respGet.Code, 200)
		require.Equal(t, len(respEvents), 1)
	})
}
