//go:build integration

package integrationtests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/postgres"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

type APITestSuite struct {
	suite.Suite
	client     http.Client
	serverAddr string
	pg         *postgres.Postgres
}

func (s *APITestSuite) SetupSuite() {
	s.serverAddr = getServerAddr(s.T())
	s.pg = getPG(s.T())
	s.client = http.Client{Timeout: 10 * time.Second}
}

func (s *APITestSuite) TearDownSuite() {
	s.pg.Close()
}

func (s *APITestSuite) TestCreateEvent() {
	s.Run("success", func() {
		event := s.createEvent()

		data := getEventDataFromDB(s.T(), s.pg, event.ID)
		s.Require().Equal(int32(1), data.Version)
	})

	s.Run("bad request", func() {
		ep := testURL(s.serverAddr, "")
		resp, err := s.client.Post(ep, ct, emptyBodyReader(s.T()))
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *APITestSuite) TestUpdateEvent() {
	s.Run("success", func() {
		event := s.createEvent()

		req, err := http.NewRequest(
			"PUT",
			s.epDetails(event.ID),
			updateBodyReader(s.T(), event),
		)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		newEvent := parseEvent(s.T(), resp.Body)
		s.Require().NotEmpty(newEvent.ID)

		data := getEventDataFromDB(s.T(), s.pg, newEvent.ID)
		s.Require().Equal(event.Version+1, data.Version)
	})

	s.Run("bad request", func() {
		event := s.createEvent()

		req, err := http.NewRequest(
			"PUT",
			s.epDetails(event.ID),
			emptyBodyReader(s.T()),
		)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

		data := getEventDataFromDB(s.T(), s.pg, event.ID)
		s.Require().Equal(event.Version, data.Version)
	})

	s.Run("incorrect version", func() {
		event := s.createEvent()

		req, err := http.NewRequest(
			"PUT",
			s.epDetails(event.ID),
			updateBodyReader(s.T(), event),
		)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		newReq, err := http.NewRequest(
			"PUT",
			s.epDetails(event.ID),
			updateBodyReader(s.T(), event),
		)
		s.Require().NoError(err)
		newResp, err := s.client.Do(newReq)
		defer func() {
			s.Require().NoError(newResp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusConflict, newResp.StatusCode)
	})

	s.Run("not found", func() {
		event := s.createEvent()
		event.ID = uuid.New()

		req, err := http.NewRequest(
			"PUT",
			s.epDetails(event.ID),
			updateBodyReader(s.T(), event),
		)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *APITestSuite) TestDeleteEvent() {
	s.Run("success", func() {
		event := s.createEvent()

		req, err := http.NewRequest("DELETE", s.epDetails(event.ID), nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
	})

	s.Run("not found", func() {
		req, err := http.NewRequest("DELETE", s.epDetails(uuid.New()), nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *APITestSuite) TestGetEventsByPeriod() {
	s.Run("success", func() {
		to := time.Now().UTC()
		from := to.AddDate(0, 0, -1)
		params := url.Values{}
		params.Add("from", from.Format(time.RFC3339))
		params.Add("to", to.Format(time.RFC3339))

		ep := testURL(s.serverAddr, "") + "?" + params.Encode()
		req, err := http.NewRequest("GET", ep, nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var events []internal.Event
		err = json.NewDecoder(resp.Body).Decode(&events)
		s.Require().NoError(err)
	})

	s.Run("bad request", func() {
		req, err := http.NewRequest("GET", testURL(s.serverAddr, ""), nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *APITestSuite) TestGetEvent() {
	s.Run("success", func() {
		event := s.createEvent()

		req, err := http.NewRequest("GET", s.epDetails(event.ID), nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		state := parseEvent(s.T(), resp.Body)
		s.Require().Equal(event.ID, state.ID)
		s.Require().Equal(event.Version, state.Version)
	})

	s.Run("not found", func() {
		req, err := http.NewRequest("GET", s.epDetails(uuid.New()), nil)
		s.Require().NoError(err)
		resp, err := s.client.Do(req)
		defer func() {
			s.Require().NoError(resp.Body.Close())
		}()

		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *APITestSuite) createEvent() internal.Event {
	ep := testURL(s.serverAddr, "")
	resp, err := s.client.Post(ep, ct, createBodyReader(s.T()))
	defer func() {
		s.Require().NoError(resp.Body.Close())
	}()

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	event := parseEvent(s.T(), resp.Body)
	s.Require().NotEmpty(event.ID)
	return event
}

func (s *APITestSuite) epDetails(ID uuid.UUID) string {
	return testURL(s.serverAddr, ID.String())
}
