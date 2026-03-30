package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort     string
	RunUpMigration bool
	PostgresConfig
}

type PostgresConfig struct {
	DBhost     string
	DBport     string
	DBuser     string
	DBpassword string
	DBname     string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// prefere loaded variables than from file
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No config file found")
	}

	return Config{
		ServerPort:     viper.GetString("SERVER_PORT"),
		RunUpMigration: viper.GetBool("RUN_UP_MIGRATION"),

		PostgresConfig: PostgresConfig{
			DBhost:     viper.GetString("DB_HOST"),
			DBport:     viper.GetString("DB_PORT"),
			DBuser:     viper.GetString("DB_USER"),
			DBpassword: viper.GetString("DB_PASSWORD"),
			DBname:     viper.GetString("DB_NAME"),
		},
	}
}
