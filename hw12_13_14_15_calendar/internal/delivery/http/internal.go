package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type InternalHTTPHandler struct {
	log logger.AppLog
	uc  delivery.HealthCheckUCI
}

func NewInternalHTTPHandler(log logger.AppLog, uc delivery.HealthCheckUCI) InternalHTTPHandler {
	return InternalHTTPHandler{log: log, uc: uc}
}

func (h *InternalHTTPHandler) Ping(c *gin.Context) {
	var ok bool
	if err := h.uc.Ping(c); err == nil {
		ok = true
	} else {
		h.log.WithReqID(c).Error("failed ping", err)
	}

	c.JSON(http.StatusOK, gin.H{"ok": ok})
}
