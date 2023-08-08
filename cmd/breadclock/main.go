package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"syscall"
)

type Config struct {
	Env  string `mapstructure:"ENV"`
	Port int    `mapstructure:"PORT"`
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
	viper.SetEnvPrefix("bc")
	viper.AutomaticEnv()

	viper.SetDefault("ENV", "development")

	_ = viper.BindEnv("ENV")
	viper.MustBindEnv("PORT")

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

	if config.Env == "production" {
		gin.SetMode("release")
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if err := r.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		zap.S().Errorw("error occurred while running http server", err)
	}
}
