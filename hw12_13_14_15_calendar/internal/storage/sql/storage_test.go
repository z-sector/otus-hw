package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

func TestPgRepo_Create(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	e := makeEvent()
	data := dto.CreateEventDTO{
		Title:            e.Title,
		BeginTime:        e.BeginTime,
		EndTime:          e.EndTime,
		Description:      e.Description,
		UserID:           e.UserID,
		NotificationTime: e.NotificationTime,
	}

	mockPool.ExpectExec("INSERT").WithArgs(
		pgxmock.AnyArg(), e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.Version, e.NotifyStatus,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	event, err := repo.Create(ctx, data)
	require.NoError(t, err)
	require.Equal(t, e.Title, event.Title)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPgRepo_Update(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	e := makeEvent()

	mockPool.ExpectExec("UPDATE").WithArgs(
		e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.Version+1, e.NotifyStatus, e.ID.String(), e.Version,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err := repo.Update(ctx, &e)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())

	mockPool.ExpectExec("UPDATE").WithArgs(
		e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.Version+1, e.NotifyStatus, e.ID.String(), e.Version,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	err = repo.Update(ctx, &e)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPgRepo_Delete(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	ID := uuid.New()
	mockPool.ExpectExec("DELETE").WithArgs(ID.String()).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(ctx, ID)
	require.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPgRepo_GetByID(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	expected := makeEvent()
	mockPool.ExpectQuery("SELECT").WithArgs(expected.ID.String()).WillReturnRows(
		pgxmock.NewRows(
			[]string{"id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status"},
		).
			AddRow(
				expected.ID,
				expected.Title,
				expected.BeginTime,
				expected.EndTime,
				expected.Description,
				expected.UserID,
				expected.NotificationTime,
				expected.Version,
				expected.NotifyStatus,
			),
	)

	actual, err := repo.GetByID(ctx, expected.ID)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPgRepo_GetByPeriod(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	e0 := makeEvent()
	e1 := makeEvent()
	from, to := time.Now().UTC(), time.Now().UTC()

	mockPool.ExpectQuery("SELECT").WithArgs(from, to).WillReturnRows(
		pgxmock.NewRows(
			[]string{"id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status"}).
			AddRow(e0.ID, e0.Title, e0.BeginTime, e0.EndTime, e0.Description, e0.UserID, e0.NotificationTime, e0.Version, e0.NotifyStatus).
			AddRow(e1.ID, e1.Title, e1.BeginTime, e1.EndTime, e1.Description, e1.UserID, e1.NotificationTime, e1.Version, e1.NotifyStatus),
	)

	eList, err := repo.GetByPeriod(ctx, from, to)
	require.NoError(t, err)
	require.Equal(t, []internal.Event{e0, e1}, eList)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPgRepo_DeleteOldEvents(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	to := time.Now().UTC()

	mockPool.ExpectExec("DELETE").WithArgs(to).WillReturnResult(pgxmock.NewResult("DELETE", 5))
	err := repo.DeleteOldEvents(ctx, to)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestPgRepo_GetEventsForNotify(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	e := makeEvent()

	mockPool.ExpectQuery("SELECT").WithArgs(*e.NotificationTime, internal.NotSentStatus).WillReturnRows(
		pgxmock.NewRows(
			[]string{"id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status"},
		).AddRow(e.ID, e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.Version, e.NotifyStatus),
	)

	eList, err := repo.GetEventsForNotify(ctx, *e.NotificationTime)
	require.NoError(t, err)
	require.Equal(t, []internal.Event{e}, eList)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPgRepo_SetNotifyStatus(t *testing.T) {
	t.Parallel()

	mockPool, repo := getMockPoolAndRepo(t)
	defer mockPool.Close()
	ctx := context.Background()

	ID := uuid.New()
	status := internal.ProcessingStatus

	mockPool.ExpectExec("UPDATE").WithArgs(status, ID.String()).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err := repo.SetNotifyStatus(ctx, ID, status)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func getMockPoolAndRepo(t *testing.T) (pgxmock.PgxPoolIface, *PgRepo) {
	t.Helper()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		require.NoError(t, err)
	}

	return mockPool, NewPgRepo(postgres.NewPg(mockPool), logger.Nop())
}

func makeEvent() internal.Event {
	notifTime := time.Now().UTC()
	beginTime := notifTime.Add(time.Second)
	endTime := beginTime.Add(time.Second)
	return internal.Event{
		ID:               uuid.New(),
		Title:            "title",
		BeginTime:        beginTime,
		EndTime:          endTime,
		Description:      "description",
		UserID:           uuid.New(),
		NotificationTime: &notifTime,
		Version:          1,
		NotifyStatus:     internal.NotSentStatus,
	}
}
