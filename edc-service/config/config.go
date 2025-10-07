package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	REST_PORT            string `mapstructure:"REST_PORT"`
	GRPC_PORT            string `mapstructure:"GRPC_PORT"`
	INTERNAL_GRPC_SECRET string `mapstructure:"INTERNAL_GRPC_SECRET"`

	JWT_SECRET            string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRATION_SECOND string `mapstructure:"JWT_EXPIRATION_SECOND"`

	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
}

var ENV Config

func LoadConfig() error {

	// * For local development
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("REST_PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No .env file found, relying on runtime environment variables.")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// * For docker environment
	v := reflect.ValueOf(ENV)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// Get the struct tag, which is the name of the environment variable
		key := t.Field(i).Tag.Get("mapstructure")
		if key != "" {
			// Bind the key to the corresponding environment variable
			// e.g., tells Viper that the key "DB_HOST" should be read from the env var "DB_HOST"
			if err := viper.BindEnv(key); err != nil {
				return err
			}
		}
	}
	viper.AutomaticEnv()

	// Unmarshal into the ENV variable
	if err := viper.Unmarshal(&ENV); err != nil {
		return err
	}

	return nil
}
