package main

import (
	"bread-clock/api"
	"bread-clock/configs"
	"bread-clock/db"
	_ "bread-clock/docs"
	"bread-clock/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"syscall"
)

// @title		빵시계 API 서버
// @version		1.0
// @description	빵시계 API 서버입니다.

// @host		110.165.23.18
// @basePath	/api/v1

// @securityDefinitions.bearer

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
	viper.MustBindEnv("DB")
	viper.MustBindEnv("DB_USER")
	viper.MustBindEnv("DB_PASSWORD")
	_ = viper.BindEnv("MIGRATE_TABLES")

	if err := viper.Unmarshal(&configs.Conf); err != nil {
		zap.S().Fatalw("failed to unmarshal config", err)
	}
}

func main() {
	defer func() {
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/breadclock?charset=utf8&parseTime=True&loc=Local", configs.Conf.DBUser, configs.Conf.DBPassword, configs.Conf.DB)
	sql, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Fatalw("failed to connect to database", err)
	}

	if configs.Conf.MigrateTables {
		migrateTables(sql)
	}

	bakeryRepository := db.NewBakeryRepository(sql)

	if configs.Conf.Env == "production" {
		gin.SetMode("release")
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.RegisterRoutes(r, bakeryRepository)

	if err := r.Run(fmt.Sprintf(":%d", configs.Conf.Port)); err != nil {
		zap.S().Errorw("error occurred while running http server", err)
	}
}

func migrateTables(sql *gorm.DB) {
	if err := sql.AutoMigrate(&models.User{}); err != nil {
		zap.S().Fatalw("failed to migrate user table", err)
	}

	if err := sql.AutoMigrate(&models.Bakery{}); err != nil {
		zap.S().Fatalw("failed to migrate bakery table", err)
	}

	if err := sql.AutoMigrate(&models.BakeryPhoto{}); err != nil {
		zap.S().Fatalw("failed to migrate bakery photo table", err)
	}

	if err := sql.AutoMigrate(&models.Bread{}); err != nil {
		zap.S().Fatalw("failed to migrate bread table", err)
	}

	if err := sql.AutoMigrate(&models.BreadPhoto{}); err != nil {
		zap.S().Fatalw("failed to migrate bread photo table", err)
	}

	if err := sql.AutoMigrate(&models.BreadAvailability{}); err != nil {
		zap.S().Fatalw("failed to migrate bread availability table", err)
	}

	if err := sql.AutoMigrate(&models.FavoriteBakery{}); err != nil {
		zap.S().Fatalw("failed to migrate favorite bakery table", err)
	}
}
