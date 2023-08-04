package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment  string `mapstructure:"ENVIRONMENT"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
}

type RedisConfig struct {
	DB       int    `yaml:"db"`
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
