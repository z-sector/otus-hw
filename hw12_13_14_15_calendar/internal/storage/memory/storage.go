package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var (
	_ usecase.EventRepo     = &MemoryRepo{}
	_ usecase.SchedulerRepo = &MemoryRepo{}
	_ usecase.SenderRepo    = &MemoryRepo{}
)

type MemoryRepo struct {
	mu     sync.RWMutex
	log    logger.AppLog
	events map[uuid.UUID]internal.Event
}

func (m *MemoryRepo) Ping(_ context.Context) error {
	return nil
}

func NewMemoryStorage(log logger.AppLog) *MemoryRepo {
	return &MemoryRepo{
		log:    log,
		events: make(map[uuid.UUID]internal.Event),
	}
}

func (m *MemoryRepo) Create(_ context.Context, data dto.CreateEventDTO) (internal.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	event := internal.Event{
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

	m.events[event.ID] = event

	return event, nil
}

func (m *MemoryRepo) Update(_ context.Context, e *internal.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	eLast, ok := m.events[e.ID]
	if !ok {
		return internal.ErrStorageNotFound
	}

	if eLast.Version != e.Version {
		return internal.ErrStorageConflict
	}
	e.Version++

	m.events[e.ID] = *e

	return nil
}

func (m *MemoryRepo) Delete(_ context.Context, ID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.existsByID(ID) {
		return internal.ErrStorageNotFound
	}

	delete(m.events, ID)

	return nil
}

func (m *MemoryRepo) GetByID(_ context.Context, ID uuid.UUID) (internal.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if e, ok := m.events[ID]; ok {
		return e, nil
	}

	return internal.Event{}, internal.ErrStorageNotFound
}

func (m *MemoryRepo) GetByPeriod(_ context.Context, from, to time.Time) ([]internal.Event, error) {
	events := make([]internal.Event, 0)

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, e := range m.events {
		if e.BeginTime.Before(to) && !e.BeginTime.Before(from) {
			events = append(events, e)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].BeginTime.Before(events[j].BeginTime)
	})

	return events, nil
}

func (m *MemoryRepo) DeleteOldEvents(_ context.Context, to time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, e := range m.events {
		if e.EndTime.Before(to) {
			delete(m.events, id)
		}
	}
	return nil
}

func (m *MemoryRepo) GetEventsForNotify(_ context.Context, time time.Time) ([]internal.Event, error) {
	events := make([]internal.Event, 0)
	for _, e := range m.events {
		if e.NotifyStatus == internal.NotSentStatus && (e.NotificationTime.After(time) || e.NotificationTime.Equal(time)) {
			events = append(events, e)
		}
	}
	return events, nil
}

func (m *MemoryRepo) SetNotifyStatus(_ context.Context, ID uuid.UUID, status internal.RemindStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := m.events[ID]; ok {
		e.NotifyStatus = status
		m.events[ID] = e
	}
	return nil
}

func (m *MemoryRepo) SetNotSentStatus(ctx context.Context, ID uuid.UUID) error {
	return m.SetNotifyStatus(ctx, ID, internal.NotSentStatus)
}

func (m *MemoryRepo) SetProcessingNotifyStatus(ctx context.Context, ID uuid.UUID) error {
	return m.SetNotifyStatus(ctx, ID, internal.ProcessingStatus)
}

func (m *MemoryRepo) SetSentNotifyStatus(ctx context.Context, ID uuid.UUID) error {
	return m.SetNotifyStatus(ctx, ID, internal.SentStatus)
}

func (m *MemoryRepo) existsByID(ID uuid.UUID) bool {
	_, ok := m.events[ID]
	return ok
}
