package internalhttp

//nolint:depguard
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/app"
	"go.uber.org/zap"
)

type Server struct {
	server  *http.Server
	logger  *zap.Logger
	handler *AppHandler
}

func NewServer(ctx context.Context, logger *zap.Logger, app *app.App) *Server {
	handler := &AppHandler{ctx: ctx, app: app, logger: logger}
	return &Server{logger: logger, handler: handler}
}

func (s *Server) Start(ctx context.Context, config configs.HTTPConfig) error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      s.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("http server is running on address: " + s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("http server is already stopped")
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/create", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.createEvent)))
	mux.HandleFunc("/update", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.updateEvent)))
	mux.HandleFunc("/delete", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.deleteEvent)))
	mux.HandleFunc("/getDay", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.getEventsDay)))
	mux.HandleFunc("/getWeek", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.getEventsWeek)))
	mux.HandleFunc("/getMonth", s.contentTypeJSONMiddleware(s.loggingMiddleware(s.handler.getEventsMonth)))

	return mux
}
