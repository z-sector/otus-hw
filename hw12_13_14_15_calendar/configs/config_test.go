package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigEnv(t *testing.T) {
	var cfg CalendarConfig
	err := ParseConfig("./config.yaml", &cfg)
	require.NoError(t, err)

	t.Setenv(fmt.Sprintf("%s_LOGGER_JSON", envPrefix), fmt.Sprintf("%t", !cfg.Logger.JSON))
	var cfgWithEnv CalendarConfig
	err = ParseConfig("./config.yaml", &cfgWithEnv)
	require.NoError(t, err)

	require.Equal(t, !cfg.Logger.JSON, cfgWithEnv.Logger.JSON)
}
