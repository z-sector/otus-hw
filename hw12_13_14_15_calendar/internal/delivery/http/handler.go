package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/delivery/http/middleware"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func NewHandler(log logger.AppLog, eUC EventUCI, iUC HealthCheckUCI) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()

	handler.Use(middleware.RequestID())
	handler.Use(middleware.Logging(log))
	handler.Use(middleware.Recovery(log))

	internalHandler := NewInternalHTTPHandler(log, iUC)
	internalRouter := handler.Group("/healthcheck")
	internalRouter.GET("/ping", internalHandler.Ping)

	eventHandler := NewEventHTTPHandler(log, eUC)
	v1Router := handler.Group("/v1")
	v1Router.GET("/test", eventHandler.Test)

	return handler
}
