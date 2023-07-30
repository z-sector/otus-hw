//go:build integration

package integrationtests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal/dto"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

func TestNotification(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

type NotificationTestSuite struct {
	suite.Suite
	client     http.Client
	serverAddr string
	pg         *postgres.Postgres
}

func (s *NotificationTestSuite) SetupSuite() {
	s.serverAddr = getServerAddr(s.T())
	s.pg = getPG(s.T())
	s.client = http.Client{Timeout: 10 * time.Second}
}

func (s *NotificationTestSuite) TearDownSuite() {
	s.pg.Close()
}

func (s *NotificationTestSuite) TestSendNotification() {
	count := 5
	events := make([]uuid.UUID, count)

	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			e := s.createEvent()
			events[i] = e.ID
		}(i)
	}
	wg.Wait()

	for i := 0; i < count; i++ {
		ID := events[i]
		s.Require().NotEqual(uuid.Nil, ID)
	}

	s.Eventually(func() bool {
		data := getCountEventsWithSentStatus(s.T(), s.pg, events)
		return data == count
	}, time.Minute, time.Second)
}

func (s *NotificationTestSuite) createEvent() internal.Event {
	BeginTime := time.Now().UTC()
	item := dto.CreateEventDTO{
		Title:            "title",
		BeginTime:        BeginTime,
		EndTime:          BeginTime.AddDate(0, 0, 1),
		Description:      "description",
		UserID:           uuid.New(),
		NotificationTime: &BeginTime,
	}
	data, err := json.Marshal(item)
	s.Require().NoError(err)

	resp, err := s.client.Post(
		testURL(s.serverAddr, ""),
		ct,
		bytes.NewReader(data),
	)
	defer func() {
		s.Require().NoError(resp.Body.Close())
	}()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	return parseEvent(s.T(), resp.Body)
}
