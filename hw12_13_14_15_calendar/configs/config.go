package configs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const envPrefix = "APP"

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

func (c Config) String() string {
	return fmt.Sprintf("Config{Logger:%s Server:%s Storage:%s}", c.Logger, c.Server, c.Storage)
}

func NewConfig(cfgFile string) (Config, error) {
	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config %q: %w", cfgFile, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode into config struct: %w", err)
	}

	if cfg.Storage.Type == StorageMemory {
		cfg.Storage.DB = nil
	}

	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validateConfig(cfg Config) error {
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
