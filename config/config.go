package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Environment  string `mapstructure:"ENVIRONMENT"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`

	GatewayAddr string   `json:"gateway_addr"`
	Servers     []string `json:"servers"`
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
		log.Fatal(err)
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
		return
	}

	return config, nil
}
