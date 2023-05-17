package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventUCI interface {
	CreateEvent(ctx context.Context, ID, title string) error
}

type EventHTTPHandler struct {
	log logger.AppLog
	uc  EventUCI
}

func NewEventHTTPHandler(log logger.AppLog, uc EventUCI) EventHTTPHandler {
	return EventHTTPHandler{log: log, uc: uc}
}

func (h *EventHTTPHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
}
