package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppB       `mapstructure:"app"`
	Server HTTPServer `mapstructure:"server"`
	DB     DataBase   `mapstructure:"database"`
}

type AppB struct {
	SercretKey string `mapstructure:"sercetKey"`
	Env        string `mapstructure:"Env"`
	Host       string `mapstructure:"Host"`
}

type HTTPServer struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

type DataBase struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	PortDb   string `mapstructure:"portDb"`
	DbName   string `mapstructure:"dbName"`
	Time     string `mapstructure:"TimeZone"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
