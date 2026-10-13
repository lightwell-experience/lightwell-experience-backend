package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// DefaultAPIPath is the path segment under /api/<path>/v.../
// Change to "lightwell" when the service leaves the lightwell-next path.
const DefaultAPIPath = "lightwell-next"

type Configuration struct {
	Loaded  bool
	API     API     `mapstructure:"api"`
	Logging Logging `mapstructure:"logging"`
	Metrics Metrics `mapstructure:"metrics"`
}

// API holds HTTP API path settings.
type API struct {
	// Path is the /api/<path>/... segment (e.g. lightwell-next, later lightwell).
	Path string `mapstructure:"path"`
}

type Logging struct {
	Level   string `mapstructure:"level"`
	Console bool   `mapstructure:"console"`
	Color   bool   `mapstructure:"color"`
}

// Metrics configures the Prometheus scrape endpoint.
type Metrics struct {
	Path string `mapstructure:"path"`
	Port int    `mapstructure:"port"`
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
	v.SetDefault("api.path", DefaultAPIPath)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9000)
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
