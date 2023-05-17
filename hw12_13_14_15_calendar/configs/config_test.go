package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigEnv(t *testing.T) {
	cfg, err := NewConfig("./config.yaml")
	require.NoError(t, err)

	t.Setenv(fmt.Sprintf("%s_LOGGER_JSON", envPrefix), fmt.Sprintf("%t", !cfg.Logger.JSON))
	cfgWithEnv, err := NewConfig("./config.yaml")
	require.NoError(t, err)

	require.Equal(t, !cfg.Logger.JSON, cfgWithEnv.Logger.JSON)
}
