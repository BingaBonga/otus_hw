package internalhttp

//nolint:depguard
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/api"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
	app    *app.App
	ctx    context.Context
}

func NewServer(ctx context.Context, logger *zap.Logger, app *app.App) *Server {
	return &Server{ctx: ctx, logger: logger, app: app}
}

func (s *Server) Start(config configs.HTTPConfig) error {
	options := api.StdHTTPServerOptions{
		BaseRouter:  http.NewServeMux(),
		Middlewares: []api.MiddlewareFunc{s.contentTypeJSONMiddleware(), s.loggingMiddleware()},
	}
	handler := api.HandlerWithOptions(s, options)
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("http server is running on address: " + s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-s.ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("http server is already stopped")
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) CreateEvent(resp http.ResponseWriter, req *http.Request) { //nolint:dupl
	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		s.logger.Error("create event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.CreateEvent(s.ctx, &event)
	if err != nil {
		s.logger.Error("create event save failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		s.logger.Error("create event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("create event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateEvent(resp http.ResponseWriter, req *http.Request) { //nolint:dupl
	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		s.logger.Error("update event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.UpdateEvent(s.ctx, &event)
	if err != nil {
		s.logger.Error("update event save failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		s.logger.Error("update event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("update event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteEvent(resp http.ResponseWriter, req *http.Request) {
	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		s.logger.Error("delete event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if event.ID == "" {
		s.logger.Error("delete event id is required")
		http.Error(resp, "id is required", http.StatusBadRequest)
		return
	}

	err = s.app.DeleteEvent(s.ctx, event.ID)
	if err != nil {
		s.logger.Error("delete event failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		s.logger.Error("delete event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("delete event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetDayEvents(resp http.ResponseWriter, _ *http.Request, owner string) { //nolint:dupl
	if owner == "" {
		s.logger.Error("get events day owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsDay(s.ctx, owner)
	if err != nil {
		s.logger.Error("get events day failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		s.logger.Error("get events day marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("get events day response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetWeekEvents(resp http.ResponseWriter, _ *http.Request, owner string) { //nolint:dupl
	if owner == "" {
		s.logger.Error("get events week owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsWeek(s.ctx, owner)
	if err != nil {
		s.logger.Error("get events week failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		s.logger.Error("get events week marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("get events week response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetMonthEvents(resp http.ResponseWriter, _ *http.Request, owner string) { //nolint: dupl
	if owner == "" {
		s.logger.Error("get events month owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsMonth(s.ctx, owner)
	if err != nil {
		s.logger.Error("get events month failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		s.logger.Error("get events month marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		s.logger.Error("get events month response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}
