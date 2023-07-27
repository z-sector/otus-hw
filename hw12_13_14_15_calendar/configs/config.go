package configs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const envPrefix = "APP"

type CalendarConfig struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

func (c CalendarConfig) String() string {
	return fmt.Sprintf("CalendarConfig{Logger:%s Server:%s Storage:%s}", c.Logger, c.Server, c.Storage)
}

type SchedulerConfig struct {
	Schedule ScheduleConf
	Logger   LoggerConf
	MQ       MQConf
	Storage  StorageConf
}

func (c SchedulerConfig) String() string {
	return fmt.Sprintf("SchedulerConfig{Schedule:%s Logger:%s MQ:%s Storage:%s}", c.Schedule, c.Logger, c.MQ, c.Storage)
}

type SenderConfig struct {
	Logger  LoggerConf
	MQ      MQConf
	Storage StorageConf
}

func (c SenderConfig) String() string {
	return fmt.Sprintf("SenderConfig{Logger:%s MQ:%s Storage:%s}", c.Logger, c.MQ, c.Storage)
}

func ParseConfig(cfgFile string, cfg any) error {
	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config %q: %w", cfgFile, err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unable to decode into config struct: %w", err)
	}

	return validateConfig(cfg)
}

func validateConfig(cfg any) error {
	err := validate.Struct(cfg)
	if err == nil {
		return nil
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		return fmt.Errorf("failed to validate config: %w", valErrs)
	}

	return fmt.Errorf("invalid configs: %w", err)
}
