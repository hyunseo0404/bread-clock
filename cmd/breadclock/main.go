package main

import (
	"errors"
	"go.uber.org/zap"
	"log"
	"os"
	"syscall"
)

func init() {
	// initialize logger
	var logger *zap.Logger
	if os.Getenv("BC_ENV") == "production" {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	defer func() {
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Printf("failed to sync logger: %v", err)
		}
	}()
}
