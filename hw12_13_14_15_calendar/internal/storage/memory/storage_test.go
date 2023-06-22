package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

const (
	eventCount        = 3
	eventDuration     = time.Minute
	eventGap          = 10 * time.Second
	eventNotifyBefore = time.Second
)

type memoryStorageTestSuite struct {
	suite.Suite
	events []internal.Event
	repo   *MemoryRepo
	begin  time.Time
}

func (s *memoryStorageTestSuite) SetupTest() {
	s.begin = time.Now().UTC()
	s.events = make([]internal.Event, eventCount)
	t := s.begin
	for i := 0; i < eventCount; i++ {
		s.events[i] = makeTestEvent(t)
		t = t.Add(eventDuration).Add(eventGap)
	}

	s.repo = NewMemoryStorage(logger.Nop())
	s.repo.events = makeEventMap(s.events)
}

func (s *memoryStorageTestSuite) TestCreate() {
	ctx := context.Background()
	s.Run("success", func() {
		newEvent := makeTestEvent(time.Now().UTC())
		data := dto.CreateEventDTO{
			Title:            newEvent.Title,
			BeginTime:        newEvent.BeginTime,
			EndTime:          newEvent.EndTime,
			Description:      newEvent.Description,
			UserID:           newEvent.UserID,
			NotificationTime: newEvent.NotificationTime,
		}

		event, err := s.repo.Create(ctx, data)
		s.Require().NoError(err)
		s.Require().True(event.ID != uuid.Nil)

		_, ok := s.repo.events[event.ID]
		s.True(ok)
	})
}

func (s *memoryStorageTestSuite) TestUpdate() {
	ctx := context.Background()
	s.Run("success", func() {
		eventID := s.events[0].ID
		changedEvent := makeTestEvent(time.Now().UTC())
		changedEvent.ID = eventID
		changedEvent.Title = "changed"

		err := s.repo.Update(ctx, &changedEvent)
		s.Require().NoError(err)

		event, ok := s.repo.events[changedEvent.ID]
		s.Require().True(ok)
		s.Require().Equal(changedEvent, event)
	})

	s.Run("error", func() {
		newEvent := makeTestEvent(time.Now().UTC())
		err := s.repo.Update(ctx, &newEvent)
		s.Require().ErrorIs(err, internal.ErrStorageNotFound)
	})
}

func (s *memoryStorageTestSuite) TestDelete() {
	ctx := context.Background()
	for _, e := range s.events {
		err := s.repo.Delete(ctx, e.ID)
		s.Require().NoError(err)

		_, ok := s.repo.events[e.ID]
		s.False(ok)

		err = s.repo.Delete(ctx, e.ID)
		s.Require().Error(err, internal.ErrStorageNotFound)
	}
}

func (s *memoryStorageTestSuite) TestGetByID() {
	ctx := context.Background()
	s.Run("success", func() {
		for _, e := range s.events {
			event, err := s.repo.GetByID(ctx, e.ID)
			s.Require().NoError(err)
			s.Require().Equal(e, event)
		}
	})

	s.Run("error", func() {
		newEvent := makeTestEvent(time.Now().UTC())

		_, err := s.repo.GetByID(ctx, newEvent.ID)
		s.Require().ErrorIs(err, internal.ErrStorageNotFound)
	})
}

func (s *memoryStorageTestSuite) TestGetByPeriod() {
	ctx := context.Background()
	s.Run("has events", func() {
		for i := 0; i < eventCount; i++ {
			for j := i; j < eventCount; j++ {
				from := s.events[i].BeginTime
				to := s.events[j].BeginTime

				events, err := s.repo.GetByPeriod(ctx, from, to)
				s.Require().NoError(err)
				s.Require().Equal(j-i, len(events))
			}
		}

		from := s.events[0].BeginTime
		to := s.events[eventCount-1].BeginTime.Add(time.Second)
		events, err := s.repo.GetByPeriod(ctx, from, to)
		s.Require().NoError(err)
		s.Require().Equal(eventCount, len(events))
	})

	s.Run("no events", func() {
		from := s.events[eventCount-1].EndTime.Add(time.Second)
		to := from.Add(time.Second)
		events, err := s.repo.GetByPeriod(ctx, from, to)
		s.Require().NoError(err)
		s.Require().Equal(0, len(events))

		to = s.events[0].BeginTime.Add(-time.Second)
		from = to.Add(-time.Second)
		events, err = s.repo.GetByPeriod(ctx, from, to)
		s.Require().NoError(err)
		s.Require().Equal(0, len(events))
	})
}

func TestMemoryStorage(t *testing.T) {
	suite.Run(t, new(memoryStorageTestSuite))
}

func TestConcurrencyUpdate(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage(logger.Nop())

	notifTime := time.Now().UTC()
	data := dto.CreateEventDTO{
		Title:            uuid.New().String(),
		BeginTime:        time.Now().UTC(),
		EndTime:          time.Now().UTC(),
		Description:      uuid.New().String(),
		UserID:           uuid.New(),
		NotificationTime: &notifTime,
	}
	event, err := storage.Create(ctx, data)
	require.NoError(t, err)

	waitCh := make(chan struct{})
	var (
		wg       sync.WaitGroup
		successC atomic.Int32
		errorC   atomic.Int32
	)

	count := 1000
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(e internal.Event) {
			defer wg.Done()
			<-waitCh
			err := storage.Update(ctx, &e)

			if err == nil {
				successC.Add(1)
				return
			}

			if errors.Is(err, internal.ErrStorageConflict) {
				errorC.Add(1)
			}
		}(event)
	}

	close(waitCh)
	wg.Wait()
	require.Equal(t, 1, int(successC.Load()))
	require.Equal(t, count-1, int(errorC.Load()))
}

func makeTestEvent(beginTime time.Time) internal.Event {
	eventID := uuid.New()
	endTime := beginTime.Add(eventDuration)
	notificationTime := beginTime.Add(-eventNotifyBefore)
	return internal.Event{
		ID:               eventID,
		Title:            fmt.Sprintf("Event_%s", eventID),
		BeginTime:        beginTime,
		EndTime:          endTime,
		Description:      fmt.Sprintf("Description: %s", eventID),
		UserID:           uuid.New(),
		NotificationTime: &notificationTime,
		Version:          1,
	}
}

func makeEventMap(events []internal.Event) map[uuid.UUID]internal.Event {
	m := make(map[uuid.UUID]internal.Event, len(events))
	for _, event := range events {
		m[event.ID] = event
	}
	return m
}
