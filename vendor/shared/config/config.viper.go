package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel string
	Port     int
	Database struct {
		ConnectString string
	}
	Auth struct {
		Url  string
		Grpc struct {
			Port int
		}
	}
	Post struct {
		Url  string
		Grpc struct {
			Port int
		}
	}
	Supabase struct {
		Url string
		Key string
	}
	Kafka struct {
		Url      string
		Username string
		Password string
	}
}

func NewConfig() *Configuration {
	fmt.Println("init config")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("config") // path to look for the config file in
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal auth config file: %w", err))
	}
	var conf Configuration
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("parse config file auth: %w", err))
	}
	return &conf
}
