package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/requestid"
)

const headerKey = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerXRequestID := c.GetHeader(headerKey)
		if headerXRequestID == "" {
			headerXRequestID = requestid.GenerateRequestID()
		}
		c.Header(headerKey, headerXRequestID)

		ctxWithReqID := requestid.CtxSetRequestID(c.Request.Context(), headerXRequestID)
		c.Request = c.Request.WithContext(ctxWithReqID)

		c.Next()
	}
}
