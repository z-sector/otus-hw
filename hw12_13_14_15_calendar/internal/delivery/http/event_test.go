package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	memorystorage "github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

func TestCreateEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log := logger.Nop()

	memStorage := memorystorage.NewMemoryStorage(log)
	eUC := usecase.NewEventUC(log, memStorage)
	handler := NewEventHTTPHandler(log, eUC)

	beginTime := time.Now().UTC()
	endTime := beginTime.Add(time.Second)
	notifTime := beginTime.Add(-time.Second)

	t.Run("success", func(t *testing.T) {
		data := dto.CreateEventDTO{
			Title:            "title",
			BeginTime:        beginTime,
			EndTime:          endTime,
			Description:      "",
			UserID:           uuid.New(),
			NotificationTime: &notifTime,
		}

		w := httptest.NewRecorder()
		ctx := createCtx(w)
		err := mockJSONPost(ctx, data)
		require.NoError(t, err)

		handler.CreateEvent(ctx)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		testcases := []struct {
			data dto.CreateEventDTO
			code int
		}{
			{
				data: dto.CreateEventDTO{
					Title:            "",
					BeginTime:        beginTime,
					EndTime:          endTime,
					Description:      "",
					UserID:           uuid.New(),
					NotificationTime: &notifTime,
				},
				code: http.StatusBadRequest,
			},
			{
				data: dto.CreateEventDTO{
					Title:            "title",
					BeginTime:        time.Time{},
					EndTime:          endTime,
					Description:      "",
					UserID:           uuid.New(),
					NotificationTime: &notifTime,
				},
				code: http.StatusBadRequest,
			},
			{
				data: dto.CreateEventDTO{
					Title:            "title",
					BeginTime:        beginTime,
					EndTime:          time.Time{},
					Description:      "",
					UserID:           uuid.New(),
					NotificationTime: &notifTime,
				},
				code: http.StatusBadRequest,
			},
			{
				data: dto.CreateEventDTO{
					Title:            "title",
					BeginTime:        beginTime,
					EndTime:          endTime,
					Description:      "",
					UserID:           uuid.Nil,
					NotificationTime: &notifTime,
				},
				code: http.StatusBadRequest,
			},
		}

		for i, tc := range testcases {
			tc := tc
			t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
				w := httptest.NewRecorder()
				ctx := createCtx(w)
				err := mockJSONPost(ctx, tc.data)
				require.NoError(t, err)

				handler.CreateEvent(ctx)
				require.Equal(t, tc.code, w.Code)
			})
		}
	})
}

func createCtx(w *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	return ctx
}

func mockJSONPost(c *gin.Context, content interface{}) error {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	return nil
}
