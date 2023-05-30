package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

func TestPgRepo(t *testing.T) {
	ctx := context.Background()
	log := logger.Nop()

	initMockPool := func() pgxmock.PgxPoolIface {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			require.NoError(t, err)
		}
		return mockPool
	}

	createPgRepo := func(pool postgres.PgPoolI) *PgRepo {
		return NewPgRepo(postgres.NewPg(pool), log)
	}

	makeEvent := func(ID uuid.UUID) internal.Event {
		notifTime := time.Now().UTC()
		beginTime := notifTime.Add(time.Second)
		endTime := beginTime.Add(time.Second)
		return internal.Event{
			ID:               ID,
			Title:            "title",
			BeginTime:        beginTime,
			EndTime:          endTime,
			Description:      "description",
			UserID:           uuid.New(),
			NotificationTime: notifTime,
		}
	}

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		mockPool := initMockPool()
		defer mockPool.Close()

		repo := createPgRepo(mockPool)

		e := makeEvent(uuid.Nil)

		mockPool.ExpectExec("INSERT").WithArgs(
			pgxmock.AnyArg(), e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime,
		).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(ctx, &e)
		require.NoError(t, err)

		err = mockPool.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		mockPool := initMockPool()
		defer mockPool.Close()

		repo := createPgRepo(mockPool)

		e := makeEvent(uuid.New())

		mockPool.ExpectExec("UPDATE").WithArgs(
			e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.ID.String(),
		).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, e)
		require.NoError(t, err)

		err = mockPool.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		mockPool := initMockPool()
		defer mockPool.Close()

		repo := createPgRepo(mockPool)

		ID := uuid.New()
		mockPool.ExpectExec("DELETE").WithArgs(ID.String()).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, ID)
		require.NoError(t, err)

		err = mockPool.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("get by id", func(t *testing.T) {
		t.Parallel()

		mockPool := initMockPool()
		defer mockPool.Close()

		repo := createPgRepo(mockPool)

		expected := makeEvent(uuid.New())
		mockPool.ExpectQuery("SELECT").WithArgs(expected.ID.String()).WillReturnRows(
			pgxmock.NewRows([]string{"id", "title", "begin_time", "end_time", "description", "user_id", "notification_time"}).
				AddRow(
					expected.ID,
					expected.Title,
					expected.BeginTime,
					expected.EndTime,
					expected.Description,
					expected.UserID,
					expected.NotificationTime,
				),
		)

		actual, err := repo.GetByID(ctx, expected.ID)
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		err = mockPool.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("get by period", func(t *testing.T) {
		t.Parallel()

		mockPool := initMockPool()
		defer mockPool.Close()

		repo := createPgRepo(mockPool)

		e0 := makeEvent(uuid.New())
		e1 := makeEvent(uuid.New())
		from, to := time.Now().UTC(), time.Now().UTC()

		mockPool.ExpectQuery("SELECT").WithArgs(from, to).WillReturnRows(
			pgxmock.NewRows(
				[]string{"id", "title", "begin_time", "end_time", "description", "user_id", "notification_time"}).
				AddRow(e0.ID, e0.Title, e0.BeginTime, e0.EndTime, e0.Description, e0.UserID, e0.NotificationTime).
				AddRow(e1.ID, e1.Title, e1.BeginTime, e1.EndTime, e1.Description, e1.UserID, e1.NotificationTime),
		)

		eList, err := repo.GetByPeriod(ctx, from, to)
		require.NoError(t, err)
		require.Equal(t, []internal.Event{e0, e1}, eList)

		err = mockPool.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
