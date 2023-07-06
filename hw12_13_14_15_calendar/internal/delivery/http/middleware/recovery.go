package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/stack"
)

func Recovery(log logger.AppLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.WithReqID(c.Request.Context())
		recovery(l)(c)
	}
}

func recovery(log logger.AppLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				brokenPipe := checkBrokenPipe(err)

				headersToStr := getHeadersStr(c.Request)

				ev := log.Base.Error().Any("error", err)

				if brokenPipe {
					ev.Str("request", headersToStr).Msg(c.Request.URL.Path)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if log.IsDebugging() {
					ev.Str("request", headersToStr)
				}
				ev.Msg(fmt.Sprintf("[Recovery from panic]:\n%s", stack.GetStack()))

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func checkBrokenPipe(a any) bool {
	if ne, ok := a.(*net.OpError); ok {
		var se *os.SyscallError
		if errors.As(ne, &se) {
			seStr := strings.ToLower(se.Error())
			if strings.Contains(seStr, "broken pipe") ||
				strings.Contains(seStr, "connection reset by peer") {
				return true
			}
		}
	}
	return false
}

func getHeadersStr(req *http.Request) string {
	httpRequest, _ := httputil.DumpRequest(req, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}
	return strings.Join(headers, "    ")
}
