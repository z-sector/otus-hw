package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logging(log logger.AppLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		baseLog := log.WithReqID(c.Request.Context()).Base.With().
			Str("ip", c.ClientIP()).
			Time("start", start).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("proto", c.Request.Proto).
			Str("agent", c.Request.UserAgent()).
			Logger()

		if log.IsDebugging() {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)

			baseLog.Debug().Any("headers", c.Request.Header).Str("body", string(body)).Msg("request")

			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		c.Next()
		end := time.Now().UTC()
		responseLog := baseLog.With().
			Time("end", end).
			Float64("latency", end.Sub(start).Seconds()).
			Int("status", c.Writer.Status()).
			Logger()

		if log.IsDebugging() {
			if blw, ok := c.Writer.(*bodyLogWriter); ok {
				responseLog.Debug().Any("headers", c.Writer.Header()).Str("body", blw.body.String()).Msg("response")
			}
		}
		responseLog.Info().Msg("response")
	}
}
