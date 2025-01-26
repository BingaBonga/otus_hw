package internalhttp

//nolint:depguard
import (
	"context"
	"net/http"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type AppHandler struct {
	app    *app.App
	ctx    context.Context
	logger *zap.Logger
}

func (handler *AppHandler) createEvent(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Body == nil || req.Body == http.NoBody {
		http.Error(resp, "body is required", http.StatusBadRequest)
		return
	}

	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		handler.logger.Error("create event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.app.CreateEvent(handler.ctx, &event)
	if err != nil {
		handler.logger.Error("create event save failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		handler.logger.Error("create event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("create event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AppHandler) updateEvent(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Body == nil || req.Body == http.NoBody {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		handler.logger.Error("update event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.app.UpdateEvent(handler.ctx, &event)
	if err != nil {
		handler.logger.Error("update event save failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		handler.logger.Error("update event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("update event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AppHandler) deleteEvent(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Body == nil || req.Body == http.NoBody {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var event storage.Event
	err := jsoniter.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		handler.logger.Error("delete event decode failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if event.ID == "" {
		handler.logger.Error("delete event id is required")
		http.Error(resp, "id is required", http.StatusBadRequest)
		return
	}

	err = handler.app.DeleteEvent(handler.ctx, event.ID)
	if err != nil {
		handler.logger.Error("delete event failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(event)
	if err != nil {
		handler.logger.Error("delete event marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("delete event response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AppHandler) getEventsDay(resp http.ResponseWriter, req *http.Request) { //nolint:dupl
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	owner := req.URL.Query().Get("owner")
	if owner == "" {
		handler.logger.Error("get events day owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := handler.app.GetEventsDay(handler.ctx, owner)
	if err != nil {
		handler.logger.Error("get events day failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		handler.logger.Error("get events day marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("get events day response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AppHandler) getEventsWeek(resp http.ResponseWriter, req *http.Request) { //nolint:dupl
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	owner := req.URL.Query().Get("owner")
	if owner == "" {
		handler.logger.Error("get events week owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := handler.app.GetEventsWeek(handler.ctx, owner)
	if err != nil {
		handler.logger.Error("get events week failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		handler.logger.Error("get events week marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("get events week response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AppHandler) getEventsMonth(resp http.ResponseWriter, req *http.Request) { //nolint: dupl
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	owner := req.URL.Query().Get("owner")
	if owner == "" {
		handler.logger.Error("get events month owner is required")
		http.Error(resp, "owner is required", http.StatusBadRequest)
		return
	}

	events, err := handler.app.GetEventsMonth(handler.ctx, owner)
	if err != nil {
		handler.logger.Error("get events month failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := jsoniter.Marshal(events)
	if err != nil {
		handler.logger.Error("get events month marshal failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		handler.logger.Error("get events month response write failed", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}
