package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/http/response"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventHTTPHandler struct {
	log logger.AppLog
	uc  delivery.EventUCI
}

func NewEventHTTPHandler(log logger.AppLog, uc delivery.EventUCI) EventHTTPHandler {
	return EventHTTPHandler{log: log, uc: uc}
}

func (h *EventHTTPHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
}

func (h *EventHTTPHandler) CreateEvent(c *gin.Context) {
	var createEvent dto.CreateEventDTO
	if err := c.ShouldBindJSON(&createEvent); err != nil {
		response.WriteError(c, http.StatusBadRequest, err)
		return
	}

	event, err := h.uc.CreateEvent(c, createEvent)
	if err != nil {
		var errVal *internal.ValidationError
		if errors.As(err, &errVal) {
			response.WriteError(c, http.StatusBadRequest, errVal)
			return
		}

		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHTTPHandler) UpdateEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusNotFound, nil)
		return
	}

	var updateEvent dto.UpdateEventDTO
	if err := c.ShouldBindJSON(&updateEvent); err != nil {
		response.WriteError(c, http.StatusBadRequest, err)
		return
	}
	updateEvent.ID = eventID

	event, err := h.uc.UpdateEvent(c, updateEvent)
	if err != nil {
		var errVal *internal.ValidationError
		if errors.As(err, &errVal) {
			response.WriteError(c, http.StatusBadRequest, errVal)
			return
		}
		if errors.Is(err, internal.ErrStorageNotFound) {
			response.WriteError(c, http.StatusNotFound, nil)
			return
		}

		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHTTPHandler) DeleteEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusNotFound, nil)
		return
	}

	if err := h.uc.DeleteEvent(c, eventID); err != nil {
		if errors.Is(err, internal.ErrStorageNotFound) {
			response.WriteError(c, http.StatusNotFound, nil)
			return
		}

		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *EventHTTPHandler) GetByID(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusNotFound, nil)
		return
	}

	event, err := h.uc.GetByID(c, eventID)
	if err != nil {
		if errors.Is(err, internal.ErrStorageNotFound) {
			response.WriteError(c, http.StatusNotFound, nil)
			return
		}

		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

type QueryParamsPeriod struct {
	From time.Time `binding:"required"`
	To   time.Time `binding:"required"`
}

func (h *EventHTTPHandler) GetByPeriod(c *gin.Context) {
	var params QueryParamsPeriod
	if err := c.ShouldBindQuery(&params); err != nil {
		response.WriteError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.uc.GetByPeriod(c, params.From, params.To)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
