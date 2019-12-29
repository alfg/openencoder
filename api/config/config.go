package config

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/viper"
)

// C is a config instance available as a public config object.
var C Config

// Config defines the main configuration object.
type Config struct {
	Port              string `mapstructure:"server_port"`
	JWTKey            string `mapstructure:"jwt_key"`
	Keyseed           string `mapstructure:"keyseed"`
	RedisHost         string `mapstructure:"redis_host"`
	RedisPort         int    `mapstructure:"redis_port"`
	RedisMaxActive    int    `mapstructure:"redis_max_active"`
	RedisMaxIdle      int    `mapstructure:"redis_max_idle"`
	DatabaseHost      string `mapstructure:"database_host"`
	DatabasePort      int    `mapstructure:"database_port"`
	DatabaseUser      string `mapstructure:"database_user"`
	DatabasePassword  string `mapstructure:"database_password"`
	DatabaseName      string `mapstructure:"database_name"`
	WorkerNamespace   string `mapstructure:"worker_namespace"`
	WorkerJobName     string `mapstructure:"worker_job_name"`
	WorkerConcurrency uint   `mapstructure:"worker_concurrency"`
	WorkDirectory     string `mapstructure:"work_dir"`

	CloudinitRedisHost        string `mapstructure:"cloudinit_redis_host"`
	CloudinitRedisPort        int    `mapstructure:"cloudinit_redis_port"`
	CloudinitDatabaseHost     string `mapstructure:"cloudinit_database_host"`
	CloudinitDatabasePort     int    `mapstructure:"cloudinit_database_port"`
	CloudinitDatabaseUser     string `mapstructure:"cloudinit_database_user"`
	CloudinitDatabasePassword string `mapstructure:"cloudinit_database_password"`
	CloudinitDatabaseName     string `mapstructure:"cloudinit_database_name"`
	CloudinitWorkerImage      string `mapstructure:"cloudinit_worker_image"`
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

// Get gets the current config.
func Get() *Config {
	return &C
}

// Keyseed gets the keyseed in a byte array.
func Keyseed() []byte {
	ks, _ := hex.DecodeString(Get().Keyseed)
	return ks
}
