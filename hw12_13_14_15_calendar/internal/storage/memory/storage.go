package memorystorage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var _ usecase.EventRepo = (*MemoryRepo)(nil)

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

func (m *MemoryRepo) Create(_ context.Context, e *internal.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	} else if m.existsByID(e.ID) {
		return fmt.Errorf("%w with id %s", internal.ErrStorageAlreadyExists, e.ID)
	}

	m.events[e.ID] = *e

	return nil
}

func (m *MemoryRepo) Update(_ context.Context, e internal.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.existsByID(e.ID) {
		return internal.ErrStorageNotFound
	}

	m.events[e.ID] = e

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

func (m *MemoryRepo) existsByID(ID uuid.UUID) bool {
	_, ok := m.events[ID]
	return ok
}
