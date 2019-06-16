package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// C is a config instance available as a public config object.
var C Config

// Config defines the main configuration object.
type Config struct {
	Port            string
	RedisHost       string `mapstructure:"redis_host"`
	RedisPort       int    `mapstructure:"redis_port"`
	WorkerNamespace string `mapstructure:"worker_namespace"`
	S3Bucket        string `mapstructure:"s3_bucket"`
	S3Region        string `mapstructure:"s3_region"`

	Profiles []profile
}

type profile struct {
	Profile string
	Output  string
	Publish bool
	Options []string
}

// LoadConfig loads up the configuration struct.
func LoadConfig(file string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	viper.AutomaticEnv()
	// config = C{}
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

// GetFFmpegProfile finds an encoding profile by profile name.
func GetFFmpegProfile(profile string) (t *profile, err error) {
	for _, v := range C.Profiles {
		if v.Profile == profile {
			return &v, nil
		}
	}
	return nil, errors.New("No task")
}

// Get gets the current config.
func Get() *Config {
	return &C
}
