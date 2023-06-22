package dto

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/require"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/internal"
)

func TestValidateTitle(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("rand seed: %d\n", seed)
	rand.Seed(seed)

	validTitle := "title"
	var invalidTitle string
	opt := options.WithRandomStringLength(101)
	require.NoError(t, faker.FakeData(&invalidTitle, opt))

	cases := []struct {
		title string
		isErr bool
	}{
		{
			title: validTitle,
			isErr: false,
		},
		{
			title: invalidTitle,
			isErr: true,
		},
	}

	for i, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actErr := ValidateTitle(tc.title)
			if tc.isErr {
				var expectedErr *internal.ValidationError
				require.ErrorAs(t, actErr, &expectedErr)
			} else {
				require.NoError(t, actErr)
			}
		})
	}
}

func TestValidateNotifTime(t *testing.T) {
	now := time.Now()

	cases := []struct {
		notif *time.Time
		begin time.Time
		isErr bool
	}{
		{
			notif: nil,
			begin: now,
			isErr: false,
		},
		{
			notif: &now,
			begin: now.Add(time.Second),
			isErr: false,
		},
		{
			notif: &now,
			begin: now,
			isErr: false,
		},
		{
			notif: &now,
			begin: now.Add(-time.Second),
			isErr: true,
		},
	}

	for i, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actErr := ValidateNotifTime(tc.notif, tc.begin)
			if tc.isErr {
				var expectedErr *internal.ValidationError
				require.ErrorAs(t, actErr, &expectedErr)
			} else {
				require.NoError(t, actErr)
			}
		})
	}
}

func TestValidateBeginEndTime(t *testing.T) {
	now := time.Now()

	cases := []struct {
		begin time.Time
		end   time.Time
		isErr bool
	}{
		{
			begin: now,
			end:   now,
			isErr: true,
		},
		{
			begin: now,
			end:   now.Add(time.Second),
			isErr: false,
		},
		{
			begin: now,
			end:   now.Add(-time.Second),
			isErr: true,
		},
	}

	for i, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actErr := ValidateBeginEndTime(tc.begin, tc.end)
			if tc.isErr {
				var expectedErr *internal.ValidationError
				require.ErrorAs(t, actErr, &expectedErr)
			} else {
				require.NoError(t, actErr)
			}
		})
	}
}
