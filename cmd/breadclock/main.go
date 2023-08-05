package main

import (
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"syscall"
)

type Config struct {
	Env string `mapstructure:"ENV"`
}

var config Config

func init() {
	// initialize logger
	var logger *zap.Logger
	if os.Getenv("BC_ENV") == "production" {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}
	zap.ReplaceGlobals(logger)

	// initialize config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.SetEnvPrefix("bc")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatalw("failed to read in config", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		zap.S().Fatalw("failed to unmarshal config", err)
	}
}

func main() {
	defer func() {
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Printf("failed to sync logger: %v", err)
		}
	}()
}
