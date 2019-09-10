package config

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// C is a config instance available as a public config object.
var C Config

// Config defines the main configuration object.
type Config struct {
	Port                    string `mapstructure:"server_port"`
	JWTKey                  string `mapstructure:"jwt_key"`
	Keyseed                 string `mapstructure:"keyseed"`
	RedisHost               string `mapstructure:"redis_host"`
	RedisPort               int    `mapstructure:"redis_port"`
	DatabaseHost            string `mapstructure:"database_host"`
	DatabasePort            int    `mapstructure:"database_port"`
	DatabaseUser            string `mapstructure:"database_user"`
	DatabasePassword        string `mapstructure:"database_password"`
	DatabaseName            string `mapstructure:"database_name"`
	WorkerNamespace         string `mapstructure:"worker_namespace"`
	WorkerJobName           string `mapstructure:"worker_job_name"`
	WorkerConcurrency       uint   `mapstructure:"worker_concurrency"`
	AWSRegion               string `mapstructure:"aws_region"`
	AWSAccessKey            string `mapstructure:"aws_access_key"`
	AWSSecretKey            string `mapstructure:"aws_secret_key"`
	S3InboundBucket         string `mapstructure:"s3_inbound_bucket"`
	S3InboundRegion         string `mapstructure:"s3_inbound_region"`
	S3OutboundBucket        string `mapstructure:"s3_outbound_bucket"`
	S3OutboundRegion        string `mapstructure:"s3_outbound_region"`
	WorkDirectory           string `mapstructure:"work_dir"`
	SlackWebhook            string `mapstructure:"slack_webhook"`
	DigitalOceanAccessToken string `mapstructure:"digitalocean_access_token"`

	CloudinitRedisHost        string `mapstructure:"cloudinit_redis_host"`
	CloudinitRedisPort        int    `mapstructure:"cloudinit_redis_port"`
	CloudinitDatabaseHost     string `mapstructure:"cloudinit_database_host"`
	CloudinitDatabasePort     int    `mapstructure:"cloudinit_database_port"`
	CloudinitDatabaseUser     string `mapstructure:"cloudinit_database_user"`
	CloudinitDatabasePassword string `mapstructure:"cloudinit_database_password"`
	CloudinitDatabaseName     string `mapstructure:"cloudinit_database_name"`

	Profiles []profile
}

type profile struct {
	Profile string   `json:"profile"`
	Output  string   `json:"output"`
	Publish bool     `json:"publish"`
	Options []string `json:"options"`
}

// LoadConfig loads up the configuration struct.
func LoadConfig(file string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	viper.AutomaticEnv()
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

// Keyseed gets the keyseed in a byte array.
func Keyseed() []byte {
	ks, _ := hex.DecodeString(Get().Keyseed)
	return ks
}
