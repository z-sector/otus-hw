//go:build integration

package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

const ct = "application/json"

type eventDataDB struct {
	ID           uuid.UUID
	Version      int32
	NotifyStatus internal.RemindStatus
}

func getServerAddr(t *testing.T) string {
	t.Helper()

	serverAddr := os.Getenv("SERVER_ADDR")
	require.NotEmpty(t, serverAddr)
	return serverAddr
}

func getPG(t *testing.T) *postgres.Postgres {
	t.Helper()

	pgAddr := os.Getenv("PG_ADDR")
	require.NotEmpty(t, pgAddr)
	pg, err := postgres.NewPgByURL(logger.GetDefaultLog(), pgAddr)
	require.NoError(t, err)
	return pg
}

func testURL(address string, ID string) string {
	if ID != "" {
		ID = "/" + ID
	}
	return fmt.Sprintf("%s/v1/events%s", address, ID)
}

func createBodyReader(t *testing.T) io.Reader {
	t.Helper()

	BeginTime := time.Now().UTC()
	item := dto.CreateEventDTO{
		Title:            "title",
		BeginTime:        BeginTime,
		EndTime:          BeginTime.AddDate(0, 0, 1),
		Description:      "description",
		UserID:           uuid.New(),
		NotificationTime: nil,
	}
	data, err := json.Marshal(item)
	require.NoError(t, err)
	return bytes.NewReader(data)
}

func emptyBodyReader(t *testing.T) io.Reader {
	t.Helper()

	return bytes.NewReader([]byte("{}"))
}

func updateBodyReader(t *testing.T, e internal.Event) io.Reader {
	t.Helper()

	item := dto.UpdateEventDTO{
		ID:               e.ID,
		Title:            fmt.Sprintf("new %s", e.Title),
		BeginTime:        e.BeginTime,
		EndTime:          e.EndTime,
		Description:      e.Description,
		UserID:           e.UserID,
		NotificationTime: e.NotificationTime,
		LastVersion:      e.Version,
	}
	data, err := json.Marshal(item)
	require.NoError(t, err)
	return bytes.NewReader(data)
}

func parseEvent(t *testing.T, closer io.ReadCloser) internal.Event {
	t.Helper()

	require.NotEmpty(t, closer)

	var event internal.Event
	err := json.NewDecoder(closer).Decode(&event)
	require.NoError(t, err)

	return event
}

func getEventDataFromDB(t *testing.T, pg *postgres.Postgres, ID uuid.UUID) eventDataDB {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sql, args, err := pg.Builder.
		Select("id", "version", "notify_status").
		From("events").
		Where(squirrel.Eq{"id": ID}).
		ToSql()
	require.NoError(t, err)

	var item eventDataDB
	err = pg.Pool.QueryRow(ctx, sql, args...).Scan(&item.ID, &item.Version, &item.NotifyStatus)
	require.NoError(t, err)
	return item
}

func getCountEventsWithSentStatus(t *testing.T, pg *postgres.Postgres, ids []uuid.UUID) int {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sql, args, err := pg.Builder.
		Select("count(*)").
		From("events").
		Where(squirrel.Eq{"id": ids, "notify_status": internal.SentStatus}).
		ToSql()
	require.NoError(t, err)

	count := 0
	err = pg.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	require.NoError(t, err)
	return count
}
