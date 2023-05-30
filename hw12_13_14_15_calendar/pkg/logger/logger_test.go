package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testCases := []string{"debug", "info", "warn", "error", "panic", "fatal"}

		for i := range testCases {
			t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
				tc := testCases[i]
				t.Parallel()

				_, err := InitLog(strings.ToLower(tc), true)
				require.NoError(t, err)

				_, err = InitLog(strings.ToUpper(tc), true)
				require.NoError(t, err)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := InitLog("test", true)
		require.Error(t, err)
	})
}
