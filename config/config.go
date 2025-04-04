package config

import "github.com/spf13/viper"

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Controller ControllerConfig `mapstructure:"controller"`
}

type AppConfig struct {
	Port int `mapstructure:"port"`
}

type ControllerConfig struct {
	Url string `mapstructure:"url"`
}

var cfg Config

func Get() *Config {
	return &cfg
}

func Init() error {
	viper.AddConfigPath("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}
