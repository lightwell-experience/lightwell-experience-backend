package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	Loaded  bool
	Logging Logging `mapstructure:"logging"`
}

type Logging struct {
	Level   string `mapstructure:"level"`
	Console bool   `mapstructure:"console"`
	Color   bool   `mapstructure:"color"`
}

var LoadedConfig Configuration

func Get() *Configuration {
	if !LoadedConfig.Loaded {
		Load()
	}
	return &LoadedConfig
}

func readConfigFile(v *viper.Viper) {
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs/")
	v.AddConfigPath("../../configs/")

	if path, ok := os.LookupEnv("CONFIG_PATH"); ok {
		v.AddConfigPath(path)
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Logger.Warn().Msgf("config.yaml file not loaded: %s", err.Error())
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("Loaded", true)
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.console", true)
	v.SetDefault("logging.color", false)
}

func Load() {
	v := viper.New()

	readConfigFile(v)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	setDefaults(v)

	if err := v.Unmarshal(&LoadedConfig); err != nil {
		panic(fmt.Errorf("unable to decode config: %w", err))
	}
}
