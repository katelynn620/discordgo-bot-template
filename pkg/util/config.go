package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Log struct {
	Level      string `mapstructure:"LEVEL"`
	MaxSize    int    `mapstructure:"MAXSIZE"`
	MaxBackups int    `mapstructure:"MAXBACKUPS"`
	MaxAge     int    `mapstructure:"MAXAGE"`
	Dir        string `mapstructure:"DIR"`
}

type Config struct {
	Debug bool   `mapstructure:"DEBUG"`
	Log   Log    `mapstructure:"LOG"`
	Token string `mapstructure:"TOKEN"`
	DB    DB     `mapstructure:"DB"`
}

type DB struct {
	Type         string `mapstructure:"TYPE"`
	File         string `mapstructure:"FILE"`
	Host         string `mapstructure:"HOST"`
	User         string `mapstructure:"USER"`
	Password     string `mapstructure:"PASSWORD"`
	Port         string `mapstructure:"PORT"`
	Name         string `mapstructure:"NAME"`
	SSL          string `mapstructure:"SSLMODE"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
}

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	// Set default values
	viper.SetDefault("debug", false)
	viper.SetDefault("log.maxsize", "1")
	viper.SetDefault("log.maxbackups", "5")
	viper.SetDefault("log.maxage", "30")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.dir", "log")
	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.file", "db.sqlite")
}

func GetConfig() (c Config) {
	if !viper.IsSet("debug") {
		LoadConfig()
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}
	return
}
