package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

var (
	_ usecase.EventRepo     = &PgRepo{}
	_ usecase.SchedulerRepo = &PgRepo{}
	_ usecase.SenderRepo    = &PgRepo{}
)

type PgRepo struct {
	*postgres.Postgres
	log       logger.AppLog
	tableName string
}

func NewPgRepo(pg *postgres.Postgres, log logger.AppLog) *PgRepo {
	return &PgRepo{Postgres: pg, log: log, tableName: "events"}
}

func (p *PgRepo) Ping(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

func (p *PgRepo) Create(ctx context.Context, data dto.CreateEventDTO) (internal.Event, error) {
	e := internal.Event{
		ID:               uuid.New(),
		Title:            data.Title,
		BeginTime:        data.BeginTime,
		EndTime:          data.EndTime,
		Description:      data.Description,
		UserID:           data.UserID,
		NotificationTime: data.NotificationTime,
		Version:          1,
		NotifyStatus:     internal.NotSentStatus,
	}

	sql, args, err := p.Builder.
		Insert(p.tableName).
		Columns("id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status").
		Values(e.ID, e.Title, e.BeginTime, e.EndTime, e.Description, e.UserID, e.NotificationTime, e.Version, e.NotifyStatus).
		ToSql()
	if err != nil {
		return internal.Event{}, fmt.Errorf("PgRepo - Create - r.Builder: %w", err)
	}

	_, err = p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return internal.Event{}, fmt.Errorf("PgRepo - Create - r.Pool.Exec: %w", err)
	}
	return e, nil
}

func (p *PgRepo) Update(ctx context.Context, e *internal.Event) error {
	e.Version++
	sql, args, err := p.Builder.
		Update(p.tableName).
		Set("title", e.Title).
		Set("begin_time", e.BeginTime).
		Set("end_time", e.EndTime).
		Set("description", e.Description).
		Set("user_id", e.UserID).
		Set("notification_time", e.NotificationTime).
		Set("version", e.Version).
		Set("notify_status", e.NotifyStatus).
		Where(squirrel.Eq{"id": e.ID, "version": e.Version - 1}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgRepo - Update - r.Builder: %w", err)
	}

	c, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PgRepo - Update - r.Pool.Exec: %w", err)
	}
	if c.RowsAffected() != 1 {
		return internal.ErrStorageConflict
	}

	return nil
}

func (p *PgRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	sql, args, err := p.Builder.Delete(p.tableName).Where(squirrel.Eq{"id": ID}).ToSql()
	if err != nil {
		return fmt.Errorf("PgRepo - Delete - r.Builder: %w", err)
	}

	c, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PgRepo - Delete - r.Pool.Exec: %w", err)
	}
	if c.RowsAffected() != 1 {
		return internal.ErrStorageNotFound
	}

	return nil
}

func (p *PgRepo) GetByID(ctx context.Context, ID uuid.UUID) (internal.Event, error) {
	sql, args, err := p.Builder.
		Select("id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version").
		From(p.tableName).
		Where(squirrel.Eq{"id": ID}).
		ToSql()
	if err != nil {
		return internal.Event{}, fmt.Errorf("PgRepo - GetByID - r.Builder: %w", err)
	}

	var e internal.Event
	err = p.Pool.QueryRow(ctx, sql, args...).Scan(
		&e.ID,
		&e.Title,
		&e.BeginTime,
		&e.EndTime,
		&e.Description,
		&e.UserID,
		&e.NotificationTime,
		&e.Version,
		&e.NotifyStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.Event{}, internal.ErrStorageNotFound
		}
		return internal.Event{}, fmt.Errorf("PgRepo - GetByID - r.Pool.QueryRow: %w", err)
	}
	return e, nil
}

func (p *PgRepo) GetByPeriod(ctx context.Context, from, to time.Time) ([]internal.Event, error) {
	filter := squirrel.And{
		squirrel.GtOrEq{"begin_time": from},
		squirrel.Lt{"begin_time": to},
	}
	sql, args, err := p.Builder.
		Select("id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status").
		From(p.tableName).Where(filter).OrderBy("begin_time").ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgRepo - GetByPeriod - r.Builder: %w", err)
	}

	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PgRepo - GetByPeriod - r.Builder: %w", err)
	}

	entities := make([]internal.Event, 0)

	for rows.Next() {
		var e internal.Event

		err = rows.Scan(
			&e.ID,
			&e.Title,
			&e.BeginTime,
			&e.EndTime,
			&e.Description,
			&e.UserID,
			&e.NotificationTime,
			&e.Version,
			&e.NotifyStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("PgRepo - GetByPeriod - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

func (p *PgRepo) DeleteOldEvents(ctx context.Context, to time.Time) error {
	filter := squirrel.Lt{"end_time": to}
	sql, args, err := p.Builder.Delete(p.tableName).Where(filter).ToSql()
	if err != nil {
		return fmt.Errorf("PgRepo - DeleteOldEvents - r.Builder: %w", err)
	}

	_, err = p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PgRepo - DeleteOldEvents - r.Pool.Exec: %w", err)
	}

	return nil
}

func (p *PgRepo) GetEventsForNotify(ctx context.Context, time time.Time) ([]internal.Event, error) {
	filter := squirrel.And{
		squirrel.NotEq{"notification_time": nil},
		squirrel.LtOrEq{"notification_time": time},
		squirrel.Eq{"notify_status": internal.NotSentStatus},
	}
	sql, args, err := p.Builder.
		Select("id", "title", "begin_time", "end_time", "description", "user_id", "notification_time", "version", "notify_status").
		From(p.tableName).Where(filter).OrderBy("notification_time").ToSql()
	if err != nil {
		return nil, fmt.Errorf("PgRepo - GetEventsForNotify - r.Builder: %w", err)
	}

	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PgRepo - GetEventsForNotify - r.Builder: %w", err)
	}

	entities := make([]internal.Event, 0)

	for rows.Next() {
		var e internal.Event

		err = rows.Scan(
			&e.ID,
			&e.Title,
			&e.BeginTime,
			&e.EndTime,
			&e.Description,
			&e.UserID,
			&e.NotificationTime,
			&e.Version,
			&e.NotifyStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("PgRepo - GetEventsForNotify - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

func (p *PgRepo) SetNotifyStatus(ctx context.Context, ID uuid.UUID, status internal.RemindStatus) error {
	filter := squirrel.Eq{"id": ID}
	sql, args, err := p.Builder.Update(p.tableName).Set("notify_status", status).Where(filter).ToSql()
	if err != nil {
		return fmt.Errorf("PgRepo - SetNotifyStatus - r.Builder: %w", err)
	}

	c, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PgRepo - SetNotifyStatus - r.Pool.Exec: %w", err)
	}
	if c.RowsAffected() != 1 {
		return internal.ErrStorageConflict
	}

	return nil
}

func (p *PgRepo) SetNotSentStatus(ctx context.Context, ID uuid.UUID) error {
	return p.SetNotifyStatus(ctx, ID, internal.NotSentStatus)
}

func (p *PgRepo) SetProcessingNotifyStatus(ctx context.Context, ID uuid.UUID) error {
	return p.SetNotifyStatus(ctx, ID, internal.ProcessingStatus)
}

func (p *PgRepo) SetSentNotifyStatus(ctx context.Context, ID uuid.UUID) error {
	return p.SetNotifyStatus(ctx, ID, internal.SentStatus)
}

func (p *PgRepo) existsByID(ctx context.Context, ID uuid.UUID) (bool, error) { // nolint: unused
	var exists bool

	sql, args, err := p.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(p.tableName).
		Where(squirrel.Eq{"id": ID}).
		Suffix(")").
		ToSql()
	if err != nil {
		return exists, fmt.Errorf("PgRepo - existsByID - r.Builder: %w", err)
	}

	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("PgRepo - existsByID - r.Pool.Exec: %w", err)
	}
	return exists, nil
}
